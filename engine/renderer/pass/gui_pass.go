package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/quad"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type UIDescriptors struct {
	descriptor.Set
	Textures *descriptor.SamplerArray
}

type UIDescriptors2 struct {
	descriptor.Set
	Config   *descriptor.Uniform[UIConfig]
	Quads    *descriptor.Storage[widget.Quad]
	Textures *descriptor.SamplerArray
}

type UIConfig struct {
	Resolution vec2.T
	ZMax       float32
}

type GuiDrawable interface {
	object.T
	DrawUI(widget.DrawArgs, *widget.QuadBuffer)
}

type GuiPass struct {
	app  vulkan.App
	mat  []material.Instance[*UIDescriptors2]
	pass renderpass.T
	fbuf framebuffer.T
	quad quad.T
}

var _ Pass = &GuiPass{}

func NewGuiPass(app vulkan.App, target RenderTarget) *GuiPass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "GUI",
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				Image:         attachment.FromImageArray(target.Output()),
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
				Name:  OutputSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{OutputAttachment},
			},
		},
	})

	mat := material.New(app.Device(), material.Args{
		Pass:       pass,
		Shader:     app.Shaders().FetchSync(shader.NewRef("ui_quad")),
		Pointers:   vertex.ParsePointers(vertex.UI{}),
		DepthTest:  true,
		DepthWrite: true,
	}, &UIDescriptors2{
		Config: &descriptor.Uniform[UIConfig]{
			Stages: core1_0.StageVertex,
		},
		Quads: &descriptor.Storage[widget.Quad]{
			Stages: core1_0.StageVertex,
			Size:   1000,
		},
		Textures: &descriptor.SamplerArray{
			Stages: core1_0.StageFragment,
			Count:  1000,
		},
	}).InstantiateMany(app.Pool(), app.Frames())

	mesh := quad.New("gui-quad", quad.Props{Size: vec2.One})

	fbufs, err := framebuffer.New(app.Device(), app.Width(), app.Height(), pass)
	if err != nil {
		panic(err)
	}

	return &GuiPass{
		app:  app,
		mat:  mat,
		pass: pass,
		fbuf: fbufs,
		quad: mesh,
	}
}

func (p *GuiPass) Record(cmds command.Recorder, args render.Args, scene object.T) {
	mat := p.mat[args.Context.Index]

	size := vec2.NewI(args.Viewport.Width, args.Viewport.Height)
	scale := args.Viewport.Scale
	size = size.Scaled(1 / scale)

	mesh, quadExists := p.app.Meshes().Fetch(p.quad.Mesh())
	if !quadExists {
		return
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf)
		mat.Bind(cmd)
	})

	// set everything to white
	white, textureReady := p.app.Textures().Fetch(color.White)
	if textureReady {
		clear := make([]texture.T, mat.Descriptors().Textures.Count)
		for i := range clear {
			clear[i] = white
		}
		mat.Descriptors().Textures.SetRange(clear, 0)
	}

	textures := cache.NewSamplerCache(p.app.Textures(), mat.Descriptors().Textures)

	uiArgs := widget.DrawArgs{
		Time:     args.Time,
		Delta:    args.Delta,
		Commands: cmds,
		Meshes:   p.app.Meshes(),
		Textures: textures,
		Viewport: render.Screen{
			Width:  int(size.X),
			Height: int(size.Y),
			Scale:  scale,
		},
	}

	qb := widget.NewQuadBuffer(100)

	// query scene for gui managers
	guis := object.Query[GuiDrawable]().Collect(scene)
	zmax := float32(0)
	for _, gui := range guis {
		gui.DrawUI(uiArgs, qb)
	}

	// todo: collect and depth sort
	mat.Descriptors().Quads.SetRange(0, qb.Data)

	for i, quad := range qb.Data {
		index := i
		zmax = math.Max(zmax, quad.ZIndex)
		cmds.Record(func(b command.Buffer) {
			mesh.Draw(b, index)
		})
	}

	mat.Descriptors().Config.Set(UIConfig{
		Resolution: size,
		ZMax:       zmax,
	})

	cmds.Record(func(cmd command.Buffer) {
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
