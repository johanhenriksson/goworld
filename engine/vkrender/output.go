package vkrender

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_shader"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_texture"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type OutputPass struct {
	backend  vulkan.T
	meshes   cache.Meshes
	textures cache.Textures
	shader   vk_shader.T[vertex.T, Uniforms, Storage]
	geometry DeferredPass

	quad *cache.VkMesh
	tex  vk_texture.T
	pass renderpass.T
}

func NewOutputPass(backend vulkan.T, meshes cache.Meshes, textures cache.Textures, geometry DeferredPass) *OutputPass {
	p := &OutputPass{
		backend:  backend,
		meshes:   meshes,
		textures: textures,
		geometry: geometry,
	}

	quadvtx := vertex.NewTriangles("screen_quad", []vertex.T{
		{P: vec3.New(-1, -1, 0), T: vec2.New(0, 0)},
		{P: vec3.New(1, 1, 0), T: vec2.New(1, 1)},
		{P: vec3.New(-1, 1, 0), T: vec2.New(0, 1)},
		{P: vec3.New(1, -1, 0), T: vec2.New(1, 0)},
	}, []uint16{
		0, 1, 2,
		0, 3, 1,
	})

	p.quad = p.meshes.Fetch(quadvtx, nil).(*cache.VkMesh)
	p.tex = p.textures.Fetch("assets/textures/uv_checker.png")

	p.pass = renderpass.New(backend.Device(), renderpass.Args{
		Frames: backend.Frames(),
		Width:  backend.Width(),
		Height: backend.Height(),
		ColorAttachments: map[string]renderpass.ColorAttachment{
			"color": {
				Images:      backend.Swapchain().Images(),
				Format:      backend.Swapchain().SurfaceFormat(),
				Samples:     vk.SampleCount1Bit,
				LoadOp:      vk.AttachmentLoadOpClear,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutPresentSrc,
			},
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:             "output",
				ColorAttachments: []string{"color"},
			},
		},
	})

	p.shader = vk_shader.New[vertex.T, Uniforms, Storage](backend, vk_shader.Args{
		Path: "vk/output",
		Pass: p.pass,
		Attributes: shader.AttributeMap{
			"position": {
				Bind: 0,
				Type: types.Float,
			},
			"texcoord_0": {
				Bind: 1,
				Type: types.Float,
			},
		},
		Samplers: vk_shader.SamplerMap{
			"diffuse": 0,
		},
	})

	return p
}

func (p *OutputPass) Draw(args render.Args, scene object.T) {
	ctx := args.Context
	worker := ctx.Workers[0]

	p.shader.SetTexture(ctx.Index, "diffuse", p.geometry.Diffuse(ctx.Index))

	worker.Queue(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)
		cmd.CmdSetViewport(0, 0, ctx.Width, ctx.Height)
		cmd.CmdSetScissor(0, 0, ctx.Width, ctx.Height)

		p.shader.Bind(ctx.Index, cmd)
		cmd.CmdBindVertexBuffer(p.quad.Vertices, 0)
		cmd.CmdBindIndexBuffers(p.quad.Indices, 0, vk.IndexTypeUint16)

		idx := 0
		cmd.CmdDrawIndexed(p.quad.Mesh.Elements(), 1, 0, 0, idx)

		cmd.CmdEndRenderPass()
	})
	worker.Submit(command.SubmitInfo{
		Wait:   []sync.Semaphore{p.geometry.Completed()},
		Signal: []sync.Semaphore{ctx.RenderComplete},
		WaitMask: []vk.PipelineStageFlags{
			vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
		},
	})
}

func (p *OutputPass) Destroy() {
	p.pass.Destroy()
	p.shader.Destroy()
}
