package vkrender

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/types"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type OutputPass struct {
	backend  vulkan.T
	meshes   cache.MeshCache
	textures cache.TextureCache
	material material.T[*OutputDescriptors]
	geometry DeferredPass

	quad      cache.VkMesh
	desc      []material.Instance[*OutputDescriptors]
	tex       []texture.T
	pass      renderpass.T
	completed sync.Semaphore
}

type OutputDescriptors struct {
	descriptor.Set
	Output *descriptor.Sampler
}

func NewOutputPass(backend vulkan.T, meshes cache.MeshCache, textures cache.TextureCache, geometry DeferredPass) *OutputPass {
	p := &OutputPass{
		backend:   backend,
		meshes:    meshes,
		textures:  textures,
		geometry:  geometry,
		completed: sync.NewSemaphore(backend.Device()),
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

	p.quad = p.meshes.Fetch(quadvtx)

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
				Usage:       vk.ImageUsageInputAttachmentBit,
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
		backend.Device(),
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
			Pass:       p.pass,
			Pointers:   vertex.ParsePointers(vertex.T{}),
			DepthTest:  false,
			DepthWrite: false,
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
		p.quad.Draw(cmd, 0)

		cmd.CmdEndRenderPass()
	})

	worker.Submit(command.SubmitInfo{
		Signal: []sync.Semaphore{p.completed},
		Wait: []command.Wait{
			{
				Semaphore: p.geometry.Completed(),
				Mask:      vk.PipelineStageColorAttachmentOutputBit,
			},
		},
	})

	// worker.Wait()
}

func (p *OutputPass) Completed() sync.Semaphore {
	return p.completed
}

func (p *OutputPass) Destroy() {
	for _, tex := range p.tex {
		tex.Destroy()
	}
	p.pass.Destroy()
	p.material.Destroy()
	p.completed.Destroy()
}
