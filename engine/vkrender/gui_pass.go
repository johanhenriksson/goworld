package vkrender

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/cache"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/material"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/shader"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/texture"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type UIDescriptors struct {
	descriptor.Set
	Textures *descriptor.SamplerArray
}

type GuiDrawable interface {
	object.Component
	DrawUI(args widget.DrawArgs, scene object.T)
}

type GuiPass struct {
	backend  vulkan.T
	mat      material.Instance[*UIDescriptors]
	pass     renderpass.T
	prev     Pass
	meshes   cache.MeshCache
	textures cache.SamplerCache
}

func NewGuiPass(backend vulkan.T, prev Pass, meshes cache.MeshCache) *GuiPass {
	pass := renderpass.New(backend.Device(), renderpass.Args{
		Frames: backend.Frames(),
		Width:  backend.Width(),
		Height: backend.Height(),
		ColorAttachments: []renderpass.ColorAttachment{
			{
				Name:          "color",
				Images:        backend.Swapchain().Images(),
				Format:        backend.Swapchain().SurfaceFormat(),
				LoadOp:        vk.AttachmentLoadOpLoad,
				StoreOp:       vk.AttachmentStoreOpStore,
				InitialLayout: vk.ImageLayoutPresentSrc,
				FinalLayout:   vk.ImageLayoutPresentSrc,
				Blend:         renderpass.BlendMultiply,
			},
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  "output",
				Depth: false,

				ColorAttachments: []string{"color"},
			},
		},
	})

	mat := material.New(backend.Device(), material.Args{
		Pass: pass,
		Shader: shader.New(backend.Device(), "vk/ui_texture", shader.Inputs{
			"position": {
				Type:  types.Float,
				Index: 0,
			},
			"color_0": {
				Type:  types.Float,
				Index: 1,
			},
			"texcoord_0": {
				Type:  types.Float,
				Index: 2,
			},
		}, shader.Descriptors{
			"Textures": 0,
		}),
		Pointers: vertex.ParsePointers(vertex.UI{}),
		Constants: []pipeline.PushConstant{
			{
				Stages: vk.ShaderStageAll,
				Type:   widget.Constants{},
			},
		},
		DepthTest:  false,
		DepthWrite: false,
	}, &UIDescriptors{
		Textures: &descriptor.SamplerArray{
			Stages: vk.ShaderStageFragmentBit,
			Count:  16,
		},
	}).Instantiate()

	textures := cache.NewSamplerCache(backend, mat.Descriptors().Textures)

	// id zero should be white
	textures.Fetch(texture.PathRef("assets/textures/white.png"))

	return &GuiPass{
		backend:  backend,
		mat:      mat,
		pass:     pass,
		prev:     prev,
		meshes:   meshes,
		textures: textures,
	}
}

func (p *GuiPass) Draw(args render.Args, scene object.T) {
	ctx := args.Context

	size := vec2.NewI(args.Viewport.Width, args.Viewport.Height)
	scale := args.Viewport.Scale
	size = size.Scaled(1 / scale)

	// setup viewport
	proj := mat4.OrthographicVK(0, size.X, 0, size.Y, 1000, -1000)
	view := mat4.Ident()
	vp := proj.Mul(&view)

	cmds := command.NewRecorder()
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)
		p.mat.Bind(cmd)
	})

	uiArgs := widget.DrawArgs{
		Commands:  cmds,
		Meshes:    p.meshes,
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
	guis := query.New[GuiDrawable]().Collect(scene)
	for _, gui := range guis {
		// todo: collect and depth sort
		gui.DrawUI(uiArgs, scene)
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})

	worker := p.backend.Worker(ctx.Index)
	worker.Queue(cmds.Apply)
	worker.Submit(command.SubmitInfo{
		Signal: []sync.Semaphore{ctx.RenderComplete},
		Wait: []command.Wait{
			{
				Semaphore: p.prev.Completed(),
				Mask:      vk.PipelineStageFragmentShaderBit,
			},
		},
	})
}

func (p *GuiPass) Destroy() {
	p.mat.Material().Destroy()
	p.pass.Destroy()
	p.textures.Destroy()
}
