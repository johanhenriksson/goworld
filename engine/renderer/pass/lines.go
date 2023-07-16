package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	lineShape "github.com/johanhenriksson/goworld/geometry/lines"
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
	fbuf      framebuffer.Array
	materials *MaterialSorter
	meshQuery *object.Query[mesh.Component]
}

func NewLinePass(app vulkan.App, target RenderTarget) *LinePass {
	log.Println("create line pass")

	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Lines",
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				Image:         attachment.FromImageArray(target.Output()),
				LoadOp:        core1_0.AttachmentLoadOpLoad,
				StoreOp:       core1_0.AttachmentStoreOpStore,
				InitialLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			},
		},
		DepthAttachment: &attachment.Depth{
			Image:         attachment.FromImageArray(target.Depth()),
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

	fbufs, err := framebuffer.NewArray(app.Frames(), app.Device(), "lines", app.Width(), app.Height(), pass)
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
		meshQuery: object.NewQuery[mesh.Component](),
	}
}

func (p *LinePass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Context.Index])
	})

	lines := p.meshQuery.
		Reset().
		Where(isDrawLines).
		Collect(scene)
	p.materials.Draw(cmds, args, lines)

	// debug lines
	if lineShape.Debug.Count() > 0 {
		lineShape.Debug.Refresh()
		p.materials.Draw(cmds, args, []mesh.Component{lineShape.Debug})
		lineShape.Debug.Clear()
	}

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

func isDrawLines(m mesh.Component) bool {
	return m.Mode() == mesh.Lines
}
