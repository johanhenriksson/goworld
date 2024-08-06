package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui"
	"github.com/johanhenriksson/goworld/gui/quad"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

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
	app    vulkan.App
	target vulkan.Target
	mat    []*material.Instance[*GuiDescriptors]
	pass   renderpass.T
	fbuf   framebuffer.Array
	quad   quad.T

	textures []cache.SamplerCache
	quads    []*widget.QuadBuffer
	guiQuery *object.Query[gui.Manager]
}

var _ Pass = &GuiPass{}

func NewGuiPass(app vulkan.App, target vulkan.Target) *GuiPass {
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

	frames := target.Frames()
	mat := material.New(app.Device(), material.Args{
		Pass:       pass,
		Shader:     app.Shaders().Fetch(shader.NewRef("ui_quad")),
		Pointers:   vertex.ParsePointers(vertex.UI{}),
		DepthTest:  true,
		DepthWrite: true,
	}, &GuiDescriptors{
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
	}).InstantiateMany(app.Pool(), target.Frames())

	mesh := quad.New("gui-quad", quad.Props{Size: vec2.One})

	fbufs, err := framebuffer.NewArray(frames, app.Device(), "gui", target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	textures := make([]cache.SamplerCache, frames)
	quads := make([]*widget.QuadBuffer, frames)
	for i := 0; i < frames; i++ {
		textures[i] = cache.NewSamplerCache(app.Textures(), mat[i].Descriptors().Textures)
		quads[i] = widget.NewQuadBuffer(10000)
	}

	return &GuiPass{
		app:      app,
		target:   target,
		mat:      mat,
		pass:     pass,
		fbuf:     fbufs,
		quad:     mesh,
		textures: textures,
		quads:    quads,
		guiQuery: object.NewQuery[gui.Manager](),
	}
}

func (p *GuiPass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	mat := p.mat[args.Frame]

	size := vec2.NewI(args.Viewport.Width, args.Viewport.Height)
	scale := args.Viewport.Scale
	size = size.Scaled(1 / scale)

	mesh := p.app.Meshes().Fetch(p.quad.Mesh())

	textures := p.textures[args.Frame]

	uiArgs := widget.DrawArgs{
		Time:     args.Time,
		Delta:    args.Delta,
		Commands: cmds,
		Textures: textures,
		Viewport: render.Screen{
			Width:  int(size.X),
			Height: int(size.Y),
			Scale:  scale,
		},
	}

	// clear quad buffer
	qb := p.quads[args.Frame]
	qb.Reset()

	// query scene for gui managers
	guis := p.guiQuery.
		Reset().
		Collect(scene)
	for _, gui := range guis {
		gui.DrawUI(uiArgs, qb)
	}

	// update sampler cache
	textures.Flush()

	// todo: collect and depth sort

	// write quad instance data
	mat.Descriptors().Quads.SetRange(0, qb.Data)

	// find maximum z
	zmax := float32(0)
	for _, quad := range qb.Data {
		zmax = math.Max(zmax, quad.ZIndex)
	}

	// update config uniform
	mat.Descriptors().Config.Set(GuiConfig{
		Resolution: size,
		ZMax:       zmax,
	})

	// draw everything in a single batch
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Frame])
		mat.Bind(cmd)
		mesh.Bind(cmd)
		mesh.DrawInstanced(cmd, 0, len(qb.Data))
		cmd.CmdEndRenderPass()
	})
}

func (p *GuiPass) Name() string {
	return "GUI"
}

func (p *GuiPass) Destroy() {
	p.mat[0].Material().Destroy()
	p.fbuf.Destroy()
	p.pass.Destroy()
}
