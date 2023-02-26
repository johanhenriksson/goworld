package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
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

type GuiDrawable interface {
	object.T
	DrawUI(args widget.DrawArgs, scene object.T)
}

type GuiPass struct {
	app  vulkan.App
	mat  []material.Instance[*UIDescriptors]
	pass renderpass.T
	fbuf framebuffer.T
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
		Pass:     pass,
		Shader:   app.Shaders().FetchSync(shader.NewRef("ui_texture")),
		Pointers: vertex.ParsePointers(vertex.UI{}),
		Constants: []pipeline.PushConstant{
			{
				Stages: core1_0.StageAll,
				Type:   widget.Constants{},
			},
		},
		DepthTest:  true,
		DepthWrite: true,
	}, &UIDescriptors{
		Textures: &descriptor.SamplerArray{
			Stages: core1_0.StageFragment,
			Count:  2000,
		},
	}).InstantiateMany(app.Pool(), app.Frames())

	fbufs, err := framebuffer.New(app.Device(), app.Width(), app.Height(), pass)
	if err != nil {
		panic(err)
	}

	return &GuiPass{
		app:  app,
		mat:  mat,
		pass: pass,
		fbuf: fbufs,
	}
}

func (p *GuiPass) Record(cmds command.Recorder, args render.Args, scene object.T) {
	mat := p.mat[args.Context.Index]

	size := vec2.NewI(args.Viewport.Width, args.Viewport.Height)
	scale := args.Viewport.Scale
	size = size.Scaled(1 / scale)

	// setup viewport
	proj := mat4.OrthographicRZ(0, size.X, 0, size.Y, 1000, -1000)
	view := mat4.Ident()
	vp := proj.Mul(&view)

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
		Commands:  cmds,
		Meshes:    p.app.Meshes(),
		Textures:  textures,
		ViewProj:  vp,
		Transform: mat4.Ident(),
		Viewport: render.Screen{
			Width:  int(size.X),
			Height: int(size.Y),
			Scale:  scale,
		},
	}

	// query scene for gui managers
	guis := object.Query[GuiDrawable]().Collect(scene)
	for _, gui := range guis {
		// todo: collect and depth sort
		gui.DrawUI(uiArgs, scene)
	}

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
