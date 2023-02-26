package pass

import (
	"fmt"
	"log"

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
	app      vulkan.App
	material material.T[*OutputDescriptors]
	source   RenderTarget

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

func NewOutputPass(app vulkan.App, source RenderTarget) *OutputPass {
	log.Println("create output pass")
	p := &OutputPass{
		app:    app,
		source: source,
	}

	p.quad = vertex.ScreenQuad("output-pass-quad")

	p.pass = renderpass.New(app.Device(), renderpass.Args{
		Name: "Output",
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Image:       attachment.FromImageArray(app.Surfaces()),
				LoadOp:      core1_0.AttachmentLoadOpClear, // clearing avoids displaying garbage on the very first frame
				FinalLayout: khr_swapchain.ImageLayoutPresentSrc,
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
		app.Device(),
		material.Args{
			Shader:     app.Shaders().FetchSync(shader.NewRef("output")),
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

	frames := app.Frames()
	var err error
	p.fbufs, err = framebuffer.NewArray(frames, app.Device(), app.Width(), app.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.desc = p.material.InstantiateMany(app.Pool(), frames)
	p.tex = make([]texture.T, frames)
	for i := range p.tex {
		outIdx := i
		if len(p.source.Output()) == 1 {
			outIdx = 0
		}
		p.tex[i], err = texture.FromImage(app.Device(), p.source.Output()[outIdx], texture.Args{
			Key:    fmt.Sprintf("gbuffer-output-%d", i),
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
	})

	quad, meshReady := p.app.Meshes().Fetch(p.quad)
	if meshReady {
		cmds.Record(func(cmd command.Buffer) {

			p.desc[ctx.Index%len(p.desc)].Bind(cmd)
			quad.Draw(cmd, 0)
		})
	}

	cmds.Record(func(cmd command.Buffer) {
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
