package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
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
	target   vulkan.Target
	material material.T[*OutputDescriptors]
	geometry GeometryBuffer
	wait     sync.Semaphore

	quad      vertex.Mesh
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

func NewOutputPass(target vulkan.Target, pool descriptor.Pool, geometry GeometryBuffer, wait sync.Semaphore) *OutputPass {
	p := &OutputPass{
		target:    target,
		geometry:  geometry,
		wait:      wait,
		completed: sync.NewSemaphore(target.Device()),
	}

	p.quad = vertex.ScreenQuad()

	p.pass = renderpass.New(target.Device(), renderpass.Args{
		ColorAttachments: []attachment.Color{
			{
				Name:        "color",
				Allocator:   attachment.FromSwapchain(target.Swapchain()),
				Format:      target.Swapchain().SurfaceFormat(),
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
		target.Device(),
		material.Args{
			Shader:     shader.New(target.Device(), "vk/output"),
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

	frames := target.Frames()
	var err error
	p.fbufs, err = framebuffer.NewArray(frames, target.Device(), target.Width(), target.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.desc = p.material.InstantiateMany(pool, frames)
	p.tex = make([]texture.T, frames)
	for i := range p.tex {
		p.tex[i], err = texture.FromView(target.Device(), p.geometry.Output(), texture.Args{
			Filter: vk.FilterNearest,
			Wrap:   vk.SamplerAddressModeClampToEdge,
		})
		if err != nil {
			// todo: clean up
			panic(err)
		}
		p.desc[i].Descriptors().Output.Set(p.tex[i])
	}

	return p
}

func (p *OutputPass) Draw(args render.Args, scene object.T) {
	ctx := args.Context

	worker := p.target.Worker(ctx.Index)
	worker.Queue(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbufs[ctx.Index%len(p.fbufs)])

		quad := p.target.Meshes().Fetch(p.quad)
		if quad != nil {
			p.desc[ctx.Index%len(p.desc)].Bind(cmd)
			quad.Draw(cmd, 0)
		}

		cmd.CmdEndRenderPass()
	})

	worker.Submit(command.SubmitInfo{
		Signal: []sync.Semaphore{p.completed},
		Wait: []command.Wait{
			{
				Semaphore: p.wait,
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
