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
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/khr_swapchain"
)

const OutputSubpass = renderpass.Name("output")

type OutputPass struct {
	target   vulkan.Target
	material material.T[*OutputDescriptors]
	geometry GeometryBuffer

	quad  vertex.Mesh
	desc  []material.Instance[*OutputDescriptors]
	tex   []texture.T
	fbufs framebuffer.Array
	pass  renderpass.T
}

type OutputDescriptors struct {
	descriptor.Set
	Output *descriptor.Sampler
}

func NewOutputPass(target vulkan.Target, geometry GeometryBuffer) *OutputPass {
	p := &OutputPass{
		target:   target,
		geometry: geometry,
	}

	p.quad = vertex.ScreenQuad("output-pass-quad")

	p.pass = renderpass.New(target.Device(), renderpass.Args{
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Allocator:   attachment.FromImageArray(target.Surfaces()),
				Format:      target.SurfaceFormat(),
				LoadOp:      core1_0.AttachmentLoadOpClear,
				FinalLayout: khr_swapchain.ImageLayoutPresentSrc,
				Usage:       core1_0.ImageUsageInputAttachment,
			},
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:             OutputSubpass,
				ColorAttachments: []attachment.Name{OutputAttachment},
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
				Stages: core1_0.StageFragment,
			},
		})

	frames := target.Frames()
	var err error
	p.fbufs, err = framebuffer.NewArray(frames, target.Device(), target.Width(), target.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.desc = p.material.InstantiateMany(target.Pool(), frames)
	p.tex = make([]texture.T, frames)
	for i := range p.tex {
		p.tex[i], err = texture.FromView(target.Device(), p.geometry.Output(), texture.Args{
			Filter: core1_0.FilterNearest,
			Wrap:   core1_0.SamplerAddressModeClampToEdge,
		})
		if err != nil {
			// todo: clean up
			panic(err)
		}
		p.desc[i].Descriptors().Output.Set(p.tex[i])
	}

	return p
}

func (p *OutputPass) Record(cmds command.Recorder, args render.Args, scene object.T) {
	ctx := args.Context
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbufs[ctx.Index%len(p.fbufs)])

		quad := p.target.Meshes().Fetch(p.quad)
		if quad != nil {
			p.desc[ctx.Index%len(p.desc)].Bind(cmd)
			quad.Draw(cmd, 0)
		}

		cmd.CmdEndRenderPass()
	})
}

func (p *OutputPass) Name() string {
	return "Output"
}

func (p *OutputPass) Destroy() {
	for _, tex := range p.tex {
		tex.Destroy()
	}
	p.fbufs.Destroy()
	p.pass.Destroy()
	p.material.Destroy()
}
