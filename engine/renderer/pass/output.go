package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/texture"
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
	fbufs     framebuffer.Array
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

	quadvtx := vertex.ScreenQuad()
	p.quad = p.meshes.Fetch(quadvtx)

	p.pass = renderpass.New(backend.Device(), renderpass.Args{
		ColorAttachments: []attachment.Color{
			{
				Name:        "color",
				Allocator:   attachment.FromSwapchain(backend.Swapchain()),
				Format:      backend.Swapchain().SurfaceFormat(),
				LoadOp:      vk.AttachmentLoadOpClear,
				FinalLayout: vk.ImageLayoutPresentSrc,
				Usage:       vk.ImageUsageInputAttachmentBit,
			},
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:             "output",
				ColorAttachments: []attachment.Name{"color"},
			},
		},
	})

	p.material = material.New(
		backend.Device(),
		material.Args{
			Shader:     shader.New(backend.Device(), "vk/output"),
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
	var err error
	p.fbufs, err = framebuffer.NewArray(frames, backend.Device(), backend.Width(), backend.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.desc = p.material.InstantiateMany(frames)
	p.tex = make([]texture.T, frames)
	for i := range p.tex {
		p.tex[i] = texture.FromView(backend.Device(), p.geometry.Output(), texture.Args{
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
		cmd.CmdBeginRenderPass(p.pass, p.fbufs[ctx.Index%len(p.fbufs)])

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
}

func (p *OutputPass) Completed() sync.Semaphore {
	return p.completed
}

func (p *OutputPass) Destroy() {
	for _, tex := range p.tex {
		tex.Destroy()
	}
	p.fbufs.Destroy()
	p.pass.Destroy()
	p.material.Destroy()
	p.completed.Destroy()
}
