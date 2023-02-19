package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type LinePass struct {
	app       vulkan.App
	pass      renderpass.T
	fbuf      framebuffer.T
	materials *MaterialSorter
}

func NewLinePass(app vulkan.App, gbuffer GeometryBuffer) *LinePass {
	log.Println("create line pass")

	pass := renderpass.New(app.Device(), renderpass.Args{
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				Image:         attachment.FromImageArray(gbuffer.Output()),
				LoadOp:        core1_0.AttachmentLoadOpLoad,
				StoreOp:       core1_0.AttachmentStoreOpStore,
				InitialLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			},
		},
		DepthAttachment: &attachment.Depth{
			Image:         attachment.FromImageArray(gbuffer.Depth()),
			LoadOp:        core1_0.AttachmentLoadOpLoad,
			InitialLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
			FinalLayout:   core1_0.ImageLayoutDepthStencilAttachmentOptimal,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  OutputSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{OutputAttachment},
			},
		},
	})

	fbufs, err := framebuffer.New(app.Device(), app.Width(), app.Height(), pass)
	if err != nil {
		panic(err)
	}

	return &LinePass{
		app:  app,
		pass: pass,
		fbuf: fbufs,
		materials: NewMaterialSorter(app, pass,
			&material.Def{
				Shader:       "lines",
				Subpass:      OutputSubpass,
				VertexFormat: vertex.C{},
				Primitive:    vertex.Lines,
				DepthTest:    true,
			}),
	}
}

func (p *LinePass) Record(cmds command.Recorder, args render.Args, scene object.T) {
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf)
	})

	lines := object.Query[mesh.T]().Where(isDrawLines).Collect(scene)
	p.materials.Draw(cmds, args, lines)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})
}

func (p *LinePass) Name() string {
	return "Lines"
}

func (p *LinePass) Destroy() {
	p.fbuf.Destroy()
	p.pass.Destroy()
	p.materials.Destroy()
}

func isDrawLines(m mesh.T) bool {
	return m.Mode() == mesh.Lines
}
