package pass

import (
	"fmt"
	"log"

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
	"github.com/vkngwrapper/extensions/v2/khr_swapchain"
)

type OutputPass struct {
	app    engine.App
	source engine.Target

	pipeline   *pipeline.Graphics
	pipeLayout *pipeline.Layout
	descLayout *descriptor.Layout[*OutputDescriptors]

	quad  vertex.Mesh
	desc  []*OutputDescriptors
	tex   texture.Array
	fbufs framebuffer.Array
	pass  *renderpass.Renderpass
}

var _ draw.Pass = &OutputPass{}

type OutputDescriptors struct {
	descriptor.Set
	Output *descriptor.Sampler
}

func NewOutputPass(app engine.App, target engine.Target, source engine.Target) *OutputPass {
	log.Println("create output pass")
	p := &OutputPass{
		app:    app,
		source: source,
	}

	p.quad = vertex.ScreenQuad("output-pass-quad")

	p.pass = renderpass.New(app.Device(), renderpass.Args{
		Name: "Output",
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Image:       attachment.FromImageArray(target.Surfaces()),
				LoadOp:      core1_0.AttachmentLoadOpClear, // clearing avoids displaying garbage on the very first frame
				FinalLayout: khr_swapchain.ImageLayoutPresentSrc,
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

	p.descLayout = descriptor.NewLayout(app.Device(), "Output", &OutputDescriptors{
		Output: &descriptor.Sampler{
			Stages: core1_0.StageFragment,
		},
	})
	p.pipeLayout = pipeline.NewLayout(app.Device(), []descriptor.SetLayout{p.descLayout}, nil)
	p.pipeline = pipeline.New(
		app.Device(),
		pipeline.Args{
			Layout:     p.pipeLayout,
			Shader:     app.Shaders().Fetch(shader.Ref("pass/output")),
			Pass:       p.pass,
			Pointers:   vertex.ParsePointers(vertex.Vertex{}),
			DepthTest:  false,
			DepthWrite: false,
		})

	frames := target.Frames()
	var err error
	p.fbufs, err = framebuffer.NewArray(frames, app.Device(), "output", target.Width(), target.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.desc = p.descLayout.InstantiateMany(app.Pool(), frames)
	p.tex = make(texture.Array, frames)
	for i := range p.tex {
		key := fmt.Sprintf("gbuffer-output-%d", i)
		p.tex[i], err = texture.FromImage(app.Device(), key, p.source.Surfaces()[i], texture.Args{
			Filter: texture.FilterNearest,
			Wrap:   texture.WrapClamp,
		})
		if err != nil {
			// todo: clean up
			panic(err)
		}
		p.desc[i].Output.Set(p.tex[i])
	}

	return p
}

func (p *OutputPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
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

func (p *OutputPass) Name() string {
	return "Output"
}

func (p *OutputPass) Destroy() {
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
