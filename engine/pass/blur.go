package pass

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type BlurPass struct {
	app   engine.App
	input engine.Target

	quad  vertex.Mesh
	desc  []*BlurDescriptors
	tex   texture.Array
	fbufs framebuffer.Array
	pass  *renderpass.Renderpass

	pipeline   *pipeline.Pipeline
	pipeLayout *pipeline.Layout
	descLayout *descriptor.Layout[*BlurDescriptors]
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
	p.descLayout = descriptor.NewLayout(app.Device(), "Blur", desc)
	p.pipeLayout = pipeline.NewLayout(app.Device(), []descriptor.SetLayout{p.descLayout}, nil)
	p.pipeline = pipeline.New(
		app.Device(),
		pipeline.Args{
			Shader:     app.Shaders().Fetch(shader.Ref("pass/blur")),
			Layout:     p.pipeLayout,
			Pass:       p.pass,
			Pointers:   vertex.ParsePointers(vertex.Vertex{}),
			DepthTest:  false,
			DepthWrite: false,
		})

	var err error
	p.fbufs, err = framebuffer.NewArray(frames, app.Device(), "blur", output.Width(), output.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.desc = p.descLayout.InstantiateMany(app.Pool(), frames)
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
		p.desc[i].Input.Set(p.tex[i])
	}

	return p
}

func (p *BlurPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	quad := p.app.Meshes().Fetch(p.quad)

	cmds.Record(func(cmd *command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbufs[args.Frame])
		cmd.CmdBindGraphicsPipeline(p.pipeline)
		cmd.CmdBindGraphicsDescriptor(p.pipeline.Layout(), 0, p.desc[args.Frame])
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
	for _, desc := range p.desc {
		desc.Destroy()
	}
	p.fbufs.Destroy()
	p.pass.Destroy()
	p.pipeline.Destroy()
	p.pipeLayout.Destroy()
	p.descLayout.Destroy()
}
