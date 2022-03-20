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
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/material"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/shader"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/texture"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type OutputPass struct {
	backend  vulkan.T
	meshes   cache.Meshes
	textures cache.Textures
	material material.T[*OutputDescriptors]
	geometry DeferredPass

	quad *cache.VkMesh
	desc []material.Instance[*OutputDescriptors]
	tex  []texture.T
	pass renderpass.T

	shadows ShadowPass
}

type OutputDescriptors struct {
	descriptor.Set
	Output *descriptor.Sampler
}

func NewOutputPass(backend vulkan.T, meshes cache.Meshes, textures cache.Textures, geometry DeferredPass, shadows ShadowPass) *OutputPass {
	p := &OutputPass{
		backend:  backend,
		meshes:   meshes,
		textures: textures,
		geometry: geometry,
		shadows:  shadows,
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

	p.pass = renderpass.New(backend.Device(), renderpass.Args{
		Frames: backend.Frames(),
		Width:  backend.Width(),
		Height: backend.Height(),
		ColorAttachments: []renderpass.ColorAttachment{
			{
				Name:        "color",
				Images:      backend.Swapchain().Images(),
				Format:      backend.Swapchain().SurfaceFormat(),
				LoadOp:      vk.AttachmentLoadOpClear,
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

	p.material = material.New(
		backend,
		material.Args{
			Shader: shader.New(
				backend.Device(),
				"vk/output",
				shader.Inputs{
					"position": {
						Index: 0,
						Type:  types.Float,
					},
					"texcoord_0": {
						Index: 1,
						Type:  types.Float,
					},
				},
				shader.Descriptors{
					"Output": 0,
				},
			),
			Pass:     p.pass,
			Pointers: vertex.ParsePointers(vertex.T{}),
		},
		&OutputDescriptors{
			Output: &descriptor.Sampler{
				Stages: vk.ShaderStageFragmentBit,
			},
		})

	frames := backend.Frames()
	p.desc = p.material.InstantiateMany(frames)
	p.tex = make([]texture.T, frames)
	for i := range p.tex {
		p.tex[i] = texture.FromView(backend.Device(), p.geometry.Output(i), texture.Args{
			Filter: vk.FilterNearest,
			Wrap:   vk.SamplerAddressModeClampToEdge,
		})
		// p.tex[i] = texture.FromView(backend.Device(), p.shadows.Shadowmap(), texture.Args{
		// 	Filter: vk.FilterNearest,
		// 	Wrap:   vk.SamplerAddressModeClampToEdge,
		// })
		p.desc[i].Descriptors().Output.Set(p.tex[i])
	}

	return p
}

func (p *OutputPass) Draw(args render.Args, scene object.T) {
	ctx := args.Context

	worker := p.backend.Worker(ctx.Index)
	worker.Queue(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, ctx.Index)

		p.desc[ctx.Index%len(p.desc)].Bind(cmd)

		cmd.CmdBindVertexBuffer(p.quad.Vertices, 0)
		cmd.CmdBindIndexBuffers(p.quad.Indices, 0, vk.IndexTypeUint16)
		cmd.CmdDrawIndexed(p.quad.Mesh.Elements(), 1, 0, 0, 0)

		cmd.CmdEndRenderPass()
	})

	worker.Submit(command.SubmitInfo{
		Signal: []sync.Semaphore{ctx.RenderComplete},
		Wait: []command.Wait{
			{
				Semaphore: p.geometry.Completed(),
				Mask:      vk.PipelineStageColorAttachmentOutputBit,
			},
		},
	})

	// worker.Wait()
}

func (p *OutputPass) Destroy() {
	for _, tex := range p.tex {
		tex.Destroy()
	}
	p.pass.Destroy()
	p.material.Destroy()
}
