package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
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

	vk "github.com/vulkan-go/vulkan"
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
	target   vulkan.Target
	mat      material.Instance[*UIDescriptors]
	pass     renderpass.T
	textures cache.SamplerCache
	fbufs    framebuffer.Array
}

var _ Pass = &GuiPass{}

func NewGuiPass(target vulkan.Target) *GuiPass {
	pass := renderpass.New(target.Device(), renderpass.Args{
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				Allocator:     attachment.FromImageArray(target.Surfaces()),
				Format:        target.SurfaceFormat(),
				LoadOp:        vk.AttachmentLoadOpLoad,
				StoreOp:       vk.AttachmentStoreOpStore,
				InitialLayout: vk.ImageLayoutPresentSrc,
				FinalLayout:   vk.ImageLayoutPresentSrc,
				Blend:         attachment.BlendMix,
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:         vk.AttachmentLoadOpClear,
			StoreOp:        vk.AttachmentStoreOpDontCare,
			StencilLoadOp:  vk.AttachmentLoadOpClear,
			StencilStoreOp: vk.AttachmentStoreOpDontCare,
			FinalLayout:    vk.ImageLayoutShaderReadOnlyOptimal,
			Usage:          vk.ImageUsageInputAttachmentBit,
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

	mat := material.New(target.Device(), material.Args{
		Pass:     pass,
		Shader:   shader.New(target.Device(), "vk/ui_texture"),
		Pointers: vertex.ParsePointers(vertex.UI{}),
		Constants: []pipeline.PushConstant{
			{
				Stages: vk.ShaderStageAll,
				Type:   widget.Constants{},
			},
		},
		DepthTest:  true,
		DepthWrite: true,
	}, &UIDescriptors{
		Textures: &descriptor.SamplerArray{
			Stages: vk.ShaderStageFragmentBit,
			Count:  2000,
		},
	}).Instantiate(target.Pool())

	fbufs, err := framebuffer.NewArray(target.Frames(), target.Device(), target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	textures := cache.NewSamplerCache(target.Device(), target.Transferer(), mat.Descriptors().Textures)
	textures.Fetch(texture.PathRef("textures/white.png")) // warmup texture

	return &GuiPass{
		target:   target,
		mat:      mat,
		pass:     pass,
		textures: textures,
		fbufs:    fbufs,
	}
}

func (p *GuiPass) Record(cmds command.Recorder, args render.Args, scene object.T) {
	p.textures.Tick()

	// texture id zero should be white
	p.textures.Fetch(texture.PathRef("textures/white.png"))

	size := vec2.NewI(args.Viewport.Width, args.Viewport.Height)
	scale := args.Viewport.Scale
	size = size.Scaled(1 / scale)

	// setup viewport
	proj := mat4.OrthographicVK(0, size.X, 0, size.Y, 1000, -1000)
	view := mat4.Ident()
	vp := proj.Mul(&view)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbufs[args.Context.Index])
		p.mat.Bind(cmd)
	})

	uiArgs := widget.DrawArgs{
		Commands:  cmds,
		Meshes:    p.target.Meshes(),
		Textures:  p.textures,
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
	p.mat.Material().Destroy()
	p.fbufs.Destroy()
	p.pass.Destroy()
	p.textures.Destroy()
}
