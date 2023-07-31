package pass

import (
	"fmt"
	"unsafe"

	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/random"
	"github.com/johanhenriksson/goworld/math/vec3"
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
	app    vulkan.App
	Target vulkan.Target
	pass   renderpass.T
	fbuf   framebuffer.Array
	mat    material.T[*AmbientOcclusionDescriptors]
	desc   []material.Instance[*AmbientOcclusionDescriptors]
	quad   vertex.Mesh

	position []texture.T
	normal   []texture.T
	kernel   [SSAOSamples]vec3.T
	noise    *HemisphereNoise
}

var _ Pass = &AmbientOcclusionPass{}

type AmbientOcclusionParams struct {
	Projection mat4.T
	Kernel     [SSAOSamples]vec3.T
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
		app: app,
	}

	// todo: optimize to single-channel texture
	p.Target, err = vulkan.NewColorTarget(app.Device(), "ssao-output", target.Width()/2, target.Height()/2, target.Frames(), target.Scale(), core1_0.FormatR8G8B8A8UnsignedNormalized)
	if err != nil {
		panic(err)
	}

	p.pass = renderpass.New(app.Device(), renderpass.Args{
		Name: "AmbientOcclusion",
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Image:       attachment.FromImageArray(p.Target.Surfaces()),
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

	p.fbuf, err = framebuffer.NewArray(target.Frames(), app.Device(), "ssao", p.Target.Width(), p.Target.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.quad = vertex.ScreenQuad("ssao-pass-quad")

	// create noise texture
	p.noise = NewHemisphereNoise(4, 4)

	// create sampler kernel
	p.kernel = [SSAOSamples]vec3.T{}
	for i := 0; i < len(p.kernel); i++ {
		sample := vec3.Random(
			vec3.New(-1, 0, -1),
			vec3.New(1, 1, 1),
		).Normalized().Scaled(random.Range(0, 1))

		// we dont want a uniform sample distribution
		// push samples closer to the origin
		scale := float32(i) / float32(SSAOSamples)
		scale = math.Lerp(0.1, 1.0, scale*scale)
		sample = sample.Scaled(scale)

		p.kernel[i] = sample
	}

	p.desc = p.mat.InstantiateMany(app.Pool(), target.Frames())
	p.position = make([]texture.T, target.Frames())
	p.normal = make([]texture.T, target.Frames())
	for i := 0; i < target.Frames(); i++ {
		posKey := fmt.Sprintf("ssao-position-%d", i)
		p.position[i], err = texture.FromImage(app.Device(), posKey, gbuffer.Position()[i], texture.Args{
			Filter: core1_0.FilterLinear,
			Wrap:   core1_0.SamplerAddressModeClampToEdge,
		})
		if err != nil {
			// todo: clean up
			panic(err)
		}
		p.desc[i].Descriptors().Position.Set(p.position[i])

		normKey := fmt.Sprintf("ssao-normal-%d", i)
		p.normal[i], err = texture.FromImage(app.Device(), normKey, gbuffer.Normal()[i], texture.Args{
			Filter: core1_0.FilterLinear,
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
	p.Target.Destroy()
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
	buffer := make([]vec3.T, 4*n.Width*n.Height)
	for i := range buffer {
		buffer[i] = vec3.Random(
			vec3.New(-1, -1, 0),
			vec3.New(1, 1, 0),
		).Normalized()
	}

	// cast to byte array
	ptr := (*byte)(unsafe.Pointer(&buffer[0]))
	bytes := unsafe.Slice(ptr, int(unsafe.Sizeof(vec3.T{}))*len(buffer))

	return &image.Data{
		Width:  n.Width,
		Height: n.Height,
		Format: image.FormatRGBA8Unorm,
		Buffer: bytes,
	}
}

func (n *HemisphereNoise) TextureArgs() texture.Args {
	return texture.Args{
		Filter: core1_0.FilterLinear,
		Wrap:   core1_0.SamplerAddressModeRepeat,
	}
}
