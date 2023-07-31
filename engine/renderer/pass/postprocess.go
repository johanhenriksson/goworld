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

type PostProcessPass struct {
	LUT texture.Ref

	app    vulkan.App
	target RenderTarget
	input  RenderTarget
	ssao   RenderTarget

	quad  vertex.Mesh
	mat   material.T[*PostProcessDescriptors]
	desc  []material.Instance[*PostProcessDescriptors]
	fbufs framebuffer.Array
	pass  renderpass.T

	inputTex []texture.T
	ssaoTex  []texture.T
}

var _ Pass = &PostProcessPass{}

type PostProcessDescriptors struct {
	descriptor.Set
	Input *descriptor.Sampler
	SSAO  *descriptor.Sampler
	LUT   *descriptor.Sampler
}

func NewPostProcessPass(app vulkan.App, input RenderTarget, ssao RenderTarget) *PostProcessPass {
	var err error
	p := &PostProcessPass{
		LUT: texture.PathRef("textures/color_grading/none.png"),

		app:   app,
		input: input,
		ssao:  ssao,
	}

	p.target, err = NewRenderTarget(app.Device(), input.Width(), input.Height(), input.Frames(),
		core1_0.FormatR8G8B8A8UnsignedNormalized, 0)
	if err != nil {
		panic(err)
	}

	p.quad = vertex.ScreenQuad("blur-pass-quad")

	p.pass = renderpass.New(app.Device(), renderpass.Args{
		Name: "PostProcess",
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Image:       attachment.FromImageArray(p.target.Output()),
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

	p.mat = material.New(
		app.Device(),
		material.Args{
			Shader:     app.Shaders().Fetch(shader.NewRef("postprocess")),
			Pass:       p.pass,
			Pointers:   vertex.ParsePointers(vertex.T{}),
			DepthTest:  false,
			DepthWrite: false,
		},
		&PostProcessDescriptors{
			Input: &descriptor.Sampler{
				Stages: core1_0.StageFragment,
			},
			SSAO: &descriptor.Sampler{
				Stages: core1_0.StageFragment,
			},
			LUT: &descriptor.Sampler{
				Stages: core1_0.StageFragment,
			},
		})

	frames := input.Frames()
	p.fbufs, err = framebuffer.NewArray(frames, app.Device(), "blur", p.target.Width(), p.target.Height(), p.pass)
	if err != nil {
		panic(err)
	}

	p.desc = p.mat.InstantiateMany(app.Pool(), frames)
	p.inputTex = make([]texture.T, frames)
	p.ssaoTex = make([]texture.T, frames)
	for i := 0; i < input.Frames(); i++ {
		inputKey := fmt.Sprintf("post-input-%d", i)
		p.inputTex[i], err = texture.FromImage(app.Device(), inputKey, p.input.Output()[i], texture.Args{
			Filter: core1_0.FilterNearest,
			Wrap:   core1_0.SamplerAddressModeClampToEdge,
		})
		if err != nil {
			// todo: clean up
			panic(err)
		}
		p.desc[i].Descriptors().Input.Set(p.inputTex[i])

		ssaoKey := fmt.Sprintf("post-ssao-%d", i)
		p.ssaoTex[i], err = texture.FromImage(app.Device(), ssaoKey, p.ssao.Output()[i], texture.Args{
			Filter: core1_0.FilterLinear,
			Wrap:   core1_0.SamplerAddressModeClampToEdge,
		})
		if err != nil {
			// todo: clean up
			panic(err)
		}
		p.desc[i].Descriptors().SSAO.Set(p.ssaoTex[i])

	}

	return p
}

func (p *PostProcessPass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	ctx := args.Context
	quad := p.app.Meshes().Fetch(p.quad)

	// refresh color lut
	lutTex := p.app.Textures().Fetch(p.LUT)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbufs[ctx.Index%len(p.fbufs)])

		desc := p.desc[ctx.Index%len(p.desc)]
		desc.Bind(cmd)
		desc.Descriptors().LUT.Set(lutTex)

		quad.Draw(cmd, 0)
		cmd.CmdEndRenderPass()
	})
}

func (p *PostProcessPass) Name() string {
	return "PostProcess"
}

func (p *PostProcessPass) Target() RenderTarget {
	return p.target
}

func (p *PostProcessPass) Destroy() {
	for _, tex := range p.inputTex {
		tex.Destroy()
	}
	for _, tex := range p.ssaoTex {
		tex.Destroy()
	}
	p.fbufs.Destroy()
	p.pass.Destroy()
	p.mat.Destroy()
}
