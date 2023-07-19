package pass

import (
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

const ForwardSubpass = renderpass.Name("forward")

type ForwardPass struct {
	gbuffer GeometryBuffer
	app     vulkan.App
	pass    renderpass.T
	fbuf    framebuffer.Array

	materials *MaterialSorter
	meshQuery *object.Query[mesh.Mesh]
}

func NewForwardPass(
	app vulkan.App,
	target RenderTarget,
	gbuffer GeometryBuffer,
) *ForwardPass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Forward",
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				LoadOp:      core1_0.AttachmentLoadOpLoad,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Blend:       attachment.BlendMultiply,

				Image: attachment.FromImageArray(target.Output()),
			},
			{
				Name:        NormalsAttachment,
				LoadOp:      core1_0.AttachmentLoadOpLoad,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,

				Image: attachment.FromImageArray(gbuffer.Normal()),
			},
			{
				Name:        PositionAttachment,
				LoadOp:      core1_0.AttachmentLoadOpLoad,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,

				Image: attachment.FromImageArray(gbuffer.Position()),
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpLoad,
			StencilLoadOp: core1_0.AttachmentLoadOpLoad,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,

			Image: attachment.FromImageArray(target.Depth()),
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  ForwardSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{OutputAttachment, NormalsAttachment, PositionAttachment},
			},
		},
	})

	fbuf, err := framebuffer.NewArray(app.Frames(), app.Device(), "forward", app.Width(), app.Height(), pass)
	if err != nil {
		panic(err)
	}

	return &ForwardPass{
		gbuffer: gbuffer,
		app:     app,
		pass:    pass,
		fbuf:    fbuf,

		materials: NewMaterialSorter(app, pass, &material.Def{
			Shader:       "color_f",
			Subpass:      ForwardSubpass,
			VertexFormat: vertex.C{},
			DepthTest:    true,
			DepthWrite:   true,
			CullMode:     vertex.CullBack,
		}),
		meshQuery: object.NewQuery[mesh.Mesh](),
	}
}

func (p *ForwardPass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	forwardMeshes := p.meshQuery.
		Reset().
		Where(isDrawForward).
		Collect(scene)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Context.Index])
	})

	p.materials.Draw(cmds, args, forwardMeshes)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})
}

func (p *ForwardPass) Name() string {
	return "Forward"
}

func (p *ForwardPass) Destroy() {
	p.fbuf.Destroy()
	p.pass.Destroy()
	p.materials.Destroy()
}

func isDrawForward(m mesh.Mesh) bool {
	return m.Mode() == mesh.Forward
}
