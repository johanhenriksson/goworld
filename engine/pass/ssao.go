package pass

import (
	"fmt"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/random"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

const SSAOSamples = 32

type AmbientOcclusionPass struct {
	app  vulkan.App
	pass renderpass.T
	fbuf framebuffer.Array
	mat  material.T[*AmbientOcclusionDescriptors]
	desc []material.Instance[*AmbientOcclusionDescriptors]
	quad vertex.Mesh

	scale    float32
	position []texture.T
	normal   []texture.T
	kernel   [SSAOSamples]vec4.T
	noise    *HemisphereNoise
}

var _ Pass = &AmbientOcclusionPass{}

type AmbientOcclusionParams struct {
	Projection mat4.T
	Kernel     [SSAOSamples]vec4.T
	Samples    int32
	Scale      float32
	Radius     float32
	Bias       float32
	Power      float32
}

type AmbientOcclusionDescriptors struct {
	descriptor.Set
	Position *descriptor.Sampler
	Normal   *descriptor.Sampler
	Noise    *descriptor.Sampler
	Params   *descriptor.Uniform[AmbientOcclusionParams]
}

func NewAmbientOcclusionPass(app vulkan.App, target vulkan.Target, gbuffer GeometryBuffer) *AmbientOcclusionPass {
	var err error
	p := &AmbientOcclusionPass{
		app:   app,
		scale: float32(gbuffer.Width()) / float32(target.Width()),
	}

	p.pass = renderpass.New(app.Device(), renderpass.Args{
		Name: "AmbientOcclusion",
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Image:       attachment.FromImageArray(target.Surfaces()),
				LoadOp:      core1_0.AttachmentLoadOpDontCare,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
			},
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
				Depth: false,

				ColorAttachments: []attachment.Name{OutputAttachment},
			},
		},
	})

	p.mat = material.New(
		app.Device(),
		material.Args{
			Shader:     app.Shaders().Fetch(shader.NewRef("ssao")),
			Pass:       p.pass,
			Pointers:   vertex.ParsePointers(vertex.T{}),
			DepthTest:  false,
			DepthWrite: false,
		},
		&AmbientOcclusionDescriptors{
			Position: &descriptor.Sampler{
				Stages: core1_0.StageFragment,
			},
			Normal: &descriptor.Sampler{
				Stages: core1_0.StageFragment,
			},
			Noise: &descriptor.Sampler{
				Stages: core1_0.StageFragment,
			},
			Params: &descriptor.Uniform[AmbientOcclusionParams]{
				Stages: core1_0.StageFragment,
			},
		})

	p.fbuf, err = framebuffer.NewArray(target.Frames(), app.Device(), "ssao", target.Width(), target.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.quad = vertex.ScreenQuad("ssao-pass-quad")

	// create noise texture
	p.noise = NewHemisphereNoise(4, 4)

	// create sampler kernel
	p.kernel = [SSAOSamples]vec4.T{}
	for i := 0; i < len(p.kernel); i++ {
		var sample vec3.T
		for {
			sample = vec3.Random(
				vec3.New(-1, 0, -1),
				vec3.New(1, 1, 1),
			)
			if sample.LengthSqr() > 1 {
				continue
			}
			sample = sample.Normalized()
			if vec3.Dot(sample, vec3.Up) < 0.5 {
				continue
			}

			sample = sample.Scaled(random.Range(0, 1))
			break
		}

		// we dont want a uniform sample distribution
		// push samples closer to the origin
		scale := float32(i) / float32(SSAOSamples)
		scale = math.Lerp(0.1, 1.0, scale*scale)
		sample = sample.Scaled(scale)

		p.kernel[i] = vec4.Extend(sample, 0)
	}

	// todo: if we shuffle the kernel, it would be ok to use fewer samples

	p.desc = p.mat.InstantiateMany(app.Pool(), target.Frames())
	p.position = make([]texture.T, target.Frames())
	p.normal = make([]texture.T, target.Frames())
	for i := 0; i < target.Frames(); i++ {
		posKey := fmt.Sprintf("ssao-position-%d", i)
		p.position[i], err = texture.FromImage(app.Device(), posKey, gbuffer.Position()[i], texture.Args{
			Filter: core1_0.FilterNearest,
			Wrap:   core1_0.SamplerAddressModeClampToEdge,
		})
		if err != nil {
			// todo: clean up
			panic(err)
		}
		p.desc[i].Descriptors().Position.Set(p.position[i])

		normKey := fmt.Sprintf("ssao-normal-%d", i)
		p.normal[i], err = texture.FromImage(app.Device(), normKey, gbuffer.Normal()[i], texture.Args{
			Filter: core1_0.FilterNearest,
			Wrap:   core1_0.SamplerAddressModeClampToEdge,
		})
		if err != nil {
			// todo: clean up
			panic(err)
		}
		p.desc[i].Descriptors().Normal.Set(p.normal[i])
	}

	return p
}

func (p *AmbientOcclusionPass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	ctx := args.Context
	quad := p.app.Meshes().Fetch(p.quad)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[ctx.Index])
		p.desc[ctx.Index].Bind(cmd)
		p.desc[ctx.Index].Descriptors().Noise.Set(p.app.Textures().Fetch(p.noise))
		p.desc[ctx.Index].Descriptors().Params.Set(AmbientOcclusionParams{
			Projection: args.Projection,
			Kernel:     p.kernel,
			Samples:    32,
			Scale:      p.scale,
			Radius:     0.7,
			Bias:       0.03,
			Power:      1.1,
		})
		quad.Draw(cmd, 0)
		cmd.CmdEndRenderPass()
	})
}

func (p *AmbientOcclusionPass) Destroy() {
	p.pass.Destroy()
	p.fbuf.Destroy()
	for i := 0; i < len(p.position); i++ {
		p.position[i].Destroy()
		p.normal[i].Destroy()
	}
	p.mat.Destroy()
}

func (p *AmbientOcclusionPass) Name() string {
	return "AmbientOcclusion"
}

type HemisphereNoise struct {
	Width  int
	Height int

	key string
}

func NewHemisphereNoise(width, height int) *HemisphereNoise {
	return &HemisphereNoise{
		key:    fmt.Sprintf("noise-hemisphere-%dx%d", width, height),
		Width:  width,
		Height: height,
	}
}

func (n *HemisphereNoise) Key() string  { return n.key }
func (n *HemisphereNoise) Version() int { return 1 }

func (n *HemisphereNoise) ImageData() *image.Data {
	buffer := make([]vec4.T, 4*n.Width*n.Height)
	for i := range buffer {
		buffer[i] = vec4.Extend(vec3.Random(
			vec3.New(-1, -1, 0),
			vec3.New(1, 1, 0),
		).Normalized(), 0)
	}

	// cast to byte array
	ptr := (*byte)(unsafe.Pointer(&buffer[0]))
	bytes := unsafe.Slice(ptr, int(unsafe.Sizeof(vec4.T{}))*len(buffer))

	return &image.Data{
		Width:  n.Width,
		Height: n.Height,
		Format: core1_0.FormatR32G32B32A32SignedFloat,
		Buffer: bytes,
	}
}

func (n *HemisphereNoise) TextureArgs() texture.Args {
	return texture.Args{
		Filter: core1_0.FilterNearest,
		Wrap:   core1_0.SamplerAddressModeRepeat,
	}
}
