package pass

import (
	"fmt"

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

type BlurPass struct {
	app      engine.App
	material *material.Material[*BlurDescriptors]
	input    engine.Target

	quad   vertex.Mesh
	layout *descriptor.Layout[*BlurDescriptors]
	desc   []*BlurDescriptors
	tex    texture.Array
	fbufs  framebuffer.Array
	pass   *renderpass.Renderpass
}

var _ draw.Pass = &BlurPass{}

type BlurDescriptors struct {
	descriptor.Set
	Input *descriptor.Sampler
}

func NewBlurPass(app engine.App, output engine.Target, input engine.Target) *BlurPass {
	p := &BlurPass{
		app:   app,
		input: input,
	}
	frames := input.Frames()

	p.quad = vertex.ScreenQuad("blur-pass-quad")

	p.pass = renderpass.New(app.Device(), renderpass.Args{
		Name: "Blur",
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Image:       attachment.FromImageArray(output.Surfaces()),
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

	desc := &BlurDescriptors{
		Input: &descriptor.Sampler{
			Stages: core1_0.StageFragment,
		},
	}
	p.layout = descriptor.NewLayout(app.Device(), "Blur", desc)
	p.material = material.New[*BlurDescriptors](
		app.Device(),
		material.Args{
			Shader:     app.Shaders().Fetch(shader.Ref("blur")),
			Pass:       p.pass,
			Pointers:   vertex.ParsePointers(vertex.T{}),
			DepthTest:  false,
			DepthWrite: false,
		},
		p.layout)

	var err error
	p.fbufs, err = framebuffer.NewArray(frames, app.Device(), "blur", output.Width(), output.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.desc = make([]*BlurDescriptors, frames)
	p.tex = make(texture.Array, frames)
	for i := range p.tex {
		key := fmt.Sprintf("blur-%d", i)
		p.tex[i], err = texture.FromImage(app.Device(), key, p.input.Surfaces()[i], texture.Args{
			Filter: texture.FilterNearest,
			Wrap:   texture.WrapClamp,
		})
		if err != nil {
			// todo: clean up
			panic(err)
		}
		p.desc[i] = p.layout.Instantiate(app.Pool())
		p.desc[i].Input.Set(p.tex[i])
	}

	return p
}

func (p *BlurPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	quad := p.app.Meshes().Fetch(p.quad)

	cmds.Record(func(cmd *command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbufs[args.Frame])
		p.material.Bind(cmd)
		cmd.CmdBindGraphicsDescriptor(p.desc[args.Frame])
		quad.Bind(cmd)
		quad.Draw(cmd, 0)
		cmd.CmdEndRenderPass()
	})
}

func (p *BlurPass) Name() string {
	return "Blur"
}

func (p *BlurPass) Destroy() {
	for _, tex := range p.tex {
		tex.Destroy()
	}
	p.fbufs.Destroy()
	p.pass.Destroy()
	p.material.Destroy()
	p.layout.Destroy()
}
