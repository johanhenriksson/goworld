package pass

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/shader"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type GuiDescriptors struct {
	descriptor.Set
	Config   *descriptor.Uniform[GuiConfig]
	Quads    *descriptor.Storage[widget.Quad]
	Textures *descriptor.SamplerArray
}

type GuiConfig struct {
	Resolution vec2.T
	ZMax       float32
}

type GuiDrawable interface {
	object.Component
	DrawUI(widget.DrawArgs, *widget.QuadBuffer)
}

type GuiPass struct {
	app    engine.App
	target engine.Target
	desc   []*GuiDescriptors
	pass   *renderpass.Renderpass
	fbuf   framebuffer.Array

	pipeline   *pipeline.Pipeline
	pipeLayout *pipeline.Layout
	descLayout *descriptor.Layout[*GuiDescriptors]

	textures cache.SamplerCache
	quads    *widget.QuadBuffer
	guiQuery *object.Query[gui.Manager]
}

var _ draw.Pass = &GuiPass{}

func NewGuiPass(app engine.App, target engine.Target) *GuiPass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "GUI",
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				Image:         attachment.FromImageArray(target.Surfaces()),
				LoadOp:        core1_0.AttachmentLoadOpLoad,
				StoreOp:       core1_0.AttachmentStoreOpStore,
				InitialLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
				Blend:         attachment.BlendMultiply,
			},
		},
		DepthAttachment: &attachment.Depth{
			Image:          attachment.NewImage("gui-depth", app.Device().GetDepthFormat(), core1_0.ImageUsageDepthStencilAttachment|core1_0.ImageUsageInputAttachment),
			LoadOp:         core1_0.AttachmentLoadOpClear,
			StoreOp:        core1_0.AttachmentStoreOpDontCare,
			StencilLoadOp:  core1_0.AttachmentLoadOpClear,
			StencilStoreOp: core1_0.AttachmentStoreOpDontCare,
			FinalLayout:    core1_0.ImageLayoutShaderReadOnlyOptimal,
			ClearDepth:     1,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{OutputAttachment},
			},
		},
		Dependencies: []renderpass.SubpassDependency{
			// fragment shader can not read the input textures until the previous pass has written to the color attachment
			{
				// For color attachment (addressing READ_AFTER_WRITE hazard)
				Src:           renderpass.ExternalSubpass,
				Dst:           MainSubpass,
				SrcStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				DstStageMask:  core1_0.PipelineStageColorAttachmentOutput,
				SrcAccessMask: core1_0.AccessColorAttachmentWrite,
				DstAccessMask: core1_0.AccessColorAttachmentRead,
				Flags:         core1_0.DependencyByRegion,
			},
			{
				// For depth attachment (addressing WRITE_AFTER_WRITE hazard)
				Src:           renderpass.ExternalSubpass,
				Dst:           MainSubpass,
				SrcStageMask:  core1_0.PipelineStageLateFragmentTests,
				DstStageMask:  core1_0.PipelineStageEarlyFragmentTests,
				SrcAccessMask: core1_0.AccessDepthStencilAttachmentWrite,
				DstAccessMask: core1_0.AccessDepthStencilAttachmentWrite | core1_0.AccessDepthStencilAttachmentRead,
				Flags:         core1_0.DependencyByRegion,
			},
		},
	})

	dlayout := descriptor.NewLayout(app.Device(), "gui", &GuiDescriptors{
		Config: &descriptor.Uniform[GuiConfig]{
			Stages: core1_0.StageAll,
		},
		Quads: &descriptor.Storage[widget.Quad]{
			Stages: core1_0.StageAll,
			Size:   10000,
		},
		Textures: &descriptor.SamplerArray{
			Stages: core1_0.StageFragment,
			Count:  1000,
		},
	})
	playout := pipeline.NewLayout(app.Device(), []descriptor.SetLayout{dlayout}, nil)

	frames := target.Frames()
	pipe := pipeline.New(app.Device(), pipeline.Args{
		Pass:       pass,
		Layout:     playout,
		Shader:     app.Shaders().Fetch(shader.Ref("pass/ui_quad")),
		DepthTest:  true,
		DepthWrite: true,
	})

	fbufs, err := framebuffer.NewArray(frames, app.Device(), "gui", target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	maxTextures := 256
	maxQuads := 10000
	desc := dlayout.InstantiateMany(app.Pool(), frames)
	textures := cache.NewSamplerCache(app.Textures(), maxTextures)
	quads := widget.NewQuadBuffer(maxQuads)

	return &GuiPass{
		app:      app,
		target:   target,
		desc:     desc,
		pass:     pass,
		fbuf:     fbufs,
		textures: textures,
		quads:    quads,
		guiQuery: object.NewQuery[gui.Manager](),

		pipeline:   pipe,
		pipeLayout: playout,
		descLayout: dlayout,
	}
}

func (p *GuiPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	desc := p.desc[args.Frame]

	size := vec2.NewI(args.Camera.Viewport.Width, args.Camera.Viewport.Height)
	scale := args.Camera.Viewport.Scale
	size = size.Scaled(1 / scale)

	uiArgs := widget.DrawArgs{
		Time:     args.Time,
		Delta:    args.Delta,
		Commands: cmds,
		Textures: p.textures,
		Viewport: draw.Viewport{
			Width:  int(size.X),
			Height: int(size.Y),
			Scale:  scale,
		},
	}

	// clear quad buffer
	p.quads.Reset()

	// query scene for gui managers
	guis := p.guiQuery.
		Reset().
		Collect(scene)
	for _, gui := range guis {
		gui.DrawUI(uiArgs, p.quads)
	}

	// update sampler cache
	p.textures.Flush(desc.Textures)

	// todo: collect and depth sort

	// write quad instance data
	desc.Quads.SetRange(0, p.quads.Data)

	// find maximum z
	zmax := float32(0)
	for _, quad := range p.quads.Data {
		zmax = math.Max(zmax, quad.ZIndex)
	}

	// update config uniform
	desc.Config.Set(GuiConfig{
		Resolution: size,
		ZMax:       zmax,
	})

	// draw everything in a single batch
	// todo: split into multiple batches to be able to use scissor test
	// todo: use draw indirect
	cmds.Record(func(cmd *command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Frame])
		cmd.CmdBindGraphicsPipeline(p.pipeline)
		cmd.CmdBindGraphicsDescriptor(p.pipeline.Layout(), 0, desc)

		cmd.CmdDraw(command.Draw{
			// the gui quad shader does not use any vertex attribute data.
			// position/uv are hardcoded in the shader, and retrieved used the vertex index.
			VertexOffset: 0,
			VertexCount:  6,

			// run an instance for each quad
			InstanceOffset: uint32(0),
			InstanceCount:  uint32(len(p.quads.Data)),
		})

		cmd.CmdEndRenderPass()
	})
}

func (p *GuiPass) Name() string {
	return "GUI"
}

func (p *GuiPass) Destroy() {
	for _, desc := range p.desc {
		desc.Destroy()
	}
	p.pipeline.Destroy()
	p.pipeLayout.Destroy()
	p.descLayout.Destroy()
	p.fbuf.Destroy()
	p.pass.Destroy()
}
