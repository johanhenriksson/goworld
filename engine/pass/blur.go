package pass

import (
	"fmt"

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
)

type BlurPass struct {
	app      vulkan.App
	target   vulkan.Target
	material material.T[*BlurDescriptors]
	input    vulkan.Target

	quad  vertex.Mesh
	desc  []material.Instance[*BlurDescriptors]
	tex   []texture.T
	fbufs framebuffer.Array
	pass  renderpass.T
}

var _ Pass = &BlurPass{}

type BlurDescriptors struct {
	descriptor.Set
	Input *descriptor.Sampler
}

func NewBlurPass(app vulkan.App, input vulkan.Target) *BlurPass {
	p := &BlurPass{
		app:   app,
		input: input,
	}
	frames := input.Frames()

	var err error
	// todo: optimize to single-channel texture
	p.target, err = vulkan.NewColorTarget(app.Device(), "blur-output", input.Width(), input.Height(), frames, input.Scale(), input.SurfaceFormat())

	p.quad = vertex.ScreenQuad("blur-pass-quad")

	p.pass = renderpass.New(app.Device(), renderpass.Args{
		Name: "Blur",
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Image:       attachment.FromImageArray(p.target.Surfaces()),
				LoadOp:      core1_0.AttachmentLoadOpDontCare,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
			},
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:             MainSubpass,
				ColorAttachments: []attachment.Name{OutputAttachment},
			},
		},
	})

	p.material = material.New(
		app.Device(),
		material.Args{
			Shader:     app.Shaders().Fetch(shader.NewRef("blur")),
			Pass:       p.pass,
			Pointers:   vertex.ParsePointers(vertex.T{}),
			DepthTest:  false,
			DepthWrite: false,
		},
		&BlurDescriptors{
			Input: &descriptor.Sampler{
				Stages: core1_0.StageFragment,
			},
		})

	p.fbufs, err = framebuffer.NewArray(frames, app.Device(), "blur", p.target.Width(), p.target.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.desc = p.material.InstantiateMany(app.Pool(), frames)
	p.tex = make([]texture.T, frames)
	for i := range p.tex {
		key := fmt.Sprintf("blur-%d", i)
		p.tex[i], err = texture.FromImage(app.Device(), key, p.input.Surfaces()[i], texture.Args{
			Filter: core1_0.FilterNearest,
			Wrap:   core1_0.SamplerAddressModeClampToEdge,
		})
		if err != nil {
			// todo: clean up
			panic(err)
		}
		p.desc[i].Descriptors().Input.Set(p.tex[i])
	}

	return p
}

func (p *BlurPass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	ctx := args.Context
	quad := p.app.Meshes().Fetch(p.quad)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbufs[ctx.Index%len(p.fbufs)])
		p.desc[ctx.Index%len(p.desc)].Bind(cmd)
		quad.Draw(cmd, 0)
		cmd.CmdEndRenderPass()
	})
}

func (p *BlurPass) Name() string {
	return "Blur"
}

func (p *BlurPass) Target() vulkan.Target {
	return p.target
}

func (p *BlurPass) Destroy() {
	for _, tex := range p.tex {
		tex.Destroy()
	}
	p.target.Destroy()
	p.fbufs.Destroy()
	p.pass.Destroy()
	p.material.Destroy()
}
