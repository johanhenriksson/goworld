package pass

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type PostProcessPass struct {
	LUT assets.Texture

	app   engine.App
	input engine.Target

	quad   vertex.Mesh
	mat    *material.Material[*PostProcessDescriptors]
	layout *descriptor.Layout[*PostProcessDescriptors]
	desc   []*PostProcessDescriptors
	fbufs  framebuffer.Array
	pass   *renderpass.Renderpass

	inputTex texture.Array
}

var _ draw.Pass = &PostProcessPass{}

type PostProcessDescriptors struct {
	descriptor.Set
	Input *descriptor.Sampler
	LUT   *descriptor.Sampler
}

func NewPostProcessPass(app engine.App, target engine.Target, input engine.Target) *PostProcessPass {
	var err error
	p := &PostProcessPass{
		LUT: texture.PathRef("textures/color_grading/none.png"),

		app:   app,
		input: input,
	}

	p.quad = vertex.ScreenQuad("blur-pass-quad")

	p.pass = renderpass.New(app.Device(), renderpass.Args{
		Name: "PostProcess",
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Image:       attachment.FromImageArray(target.Surfaces()),
				LoadOp:      core1_0.AttachmentLoadOpDontCare,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
			},
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:             MainSubpass,
				ColorAttachments: []attachment.Name{OutputAttachment},
			},
		},
		Dependencies: []renderpass.SubpassDependency{
			{
				// For color attachment operations
				Src:           renderpass.ExternalSubpass,
				Dst:           MainSubpass,
				SrcStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				DstStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				SrcAccessMask: core1_0.AccessColorAttachmentWrite,
				DstAccessMask: core1_0.AccessColorAttachmentWrite | core1_0.AccessColorAttachmentRead,
				Flags:         core1_0.DependencyByRegion,
			},
			{
				// For fragment shader reads
				Src:           renderpass.ExternalSubpass,
				Dst:           MainSubpass,
				SrcStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				DstStageMask:  core1_0.PipelineStageFragmentShader,
				SrcAccessMask: core1_0.AccessColorAttachmentWrite,
				DstAccessMask: core1_0.AccessShaderRead,
				Flags:         core1_0.DependencyByRegion,
			},
		},
	})

	p.layout = descriptor.NewLayout(app.Device(), "PostProcess", &PostProcessDescriptors{
		Input: &descriptor.Sampler{
			Stages: core1_0.StageFragment,
		},
		LUT: &descriptor.Sampler{
			Stages: core1_0.StageFragment,
		},
	})
	p.mat = material.New[*PostProcessDescriptors](
		app.Device(),
		material.Args{
			Shader:     app.Shaders().Fetch(shader.Ref("postprocess")),
			Pass:       p.pass,
			Pointers:   vertex.ParsePointers(vertex.T{}),
			DepthTest:  false,
			DepthWrite: false,
		},
		p.layout)

	frames := input.Frames()
	p.fbufs, err = framebuffer.NewArray(frames, app.Device(), "postprocess", target.Width(), target.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.desc = make([]*PostProcessDescriptors, frames)
	p.inputTex = make(texture.Array, frames)
	for i := 0; i < input.Frames(); i++ {
		inputKey := fmt.Sprintf("post-input-%d", i)
		p.inputTex[i], err = texture.FromImage(app.Device(), inputKey, p.input.Surfaces()[i], texture.Args{
			Filter: texture.FilterNearest,
			Wrap:   texture.WrapClamp,
		})
		if err != nil {
			// todo: clean up
			panic(err)
		}
		p.desc[i] = p.layout.Instantiate(app.Pool())
		p.desc[i].Input.Set(p.inputTex[i])
	}

	return p
}

func (p *PostProcessPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	quad := p.app.Meshes().Fetch(p.quad)

	// refresh color lut
	lutTex := p.app.Textures().Fetch(p.LUT)

	cmds.Record(func(cmd *command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbufs[args.Frame])

		p.mat.Bind(cmd)
		desc := p.desc[args.Frame]
		cmd.CmdBindGraphicsDescriptor(desc)
		desc.LUT.Set(lutTex)

		quad.Bind(cmd)
		quad.Draw(cmd, 0)
		cmd.CmdEndRenderPass()
	})
}

func (p *PostProcessPass) Name() string {
	return "PostProcess"
}

func (p *PostProcessPass) Destroy() {
	for _, tex := range p.inputTex {
		tex.Destroy()
	}
	p.fbufs.Destroy()
	p.pass.Destroy()
	p.mat.Destroy()
	p.layout.Destroy()
}
