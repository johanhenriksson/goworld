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
	gbuffer   GeometryBuffer
	app       vulkan.App
	pass      renderpass.T
	fbuf      framebuffer.T
	materials *MaterialSorter
}

func NewForwardPass(
	app vulkan.App,
	gbuffer GeometryBuffer,
) *ForwardPass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Format:      gbuffer.Output().Format(),
				LoadOp:      core1_0.AttachmentLoadOpLoad,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Usage:       core1_0.ImageUsageSampled,
				Blend:       attachment.BlendMultiply,

				Allocator: attachment.FromImage(gbuffer.Output()),
			},
			{
				Name:        NormalsAttachment,
				Format:      gbuffer.Normal().Format(),
				LoadOp:      core1_0.AttachmentLoadOpLoad,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Usage:       core1_0.ImageUsageInputAttachment | core1_0.ImageUsageTransferSrc,

				Allocator: attachment.FromImage(gbuffer.Normal()),
			},
			{
				Name:        PositionAttachment,
				Format:      gbuffer.Position().Format(),
				LoadOp:      core1_0.AttachmentLoadOpLoad,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Usage:       core1_0.ImageUsageInputAttachment | core1_0.ImageUsageTransferSrc,

				Allocator: attachment.FromImage(gbuffer.Position()),
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpLoad,
			StencilLoadOp: core1_0.AttachmentLoadOpLoad,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			Usage:         core1_0.ImageUsageInputAttachment,

			Allocator: attachment.FromImage(gbuffer.Depth()),
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  ForwardSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{OutputAttachment, NormalsAttachment, PositionAttachment},
			},
		},
	})

	fbuf, err := framebuffer.New(app.Device(), app.Width(), app.Height(), pass)
	if err != nil {
		panic(err)
	}

	return &ForwardPass{
		gbuffer: gbuffer,
		app:     app,
		pass:    pass,
		fbuf:    fbuf,

		materials: NewMaterialSorter(app, pass, &material.Def{
			Shader:       "vk/color_f",
			Subpass:      ForwardSubpass,
			VertexFormat: vertex.C{},
			DepthTest:    true,
			DepthWrite:   true,
			CullMode:     vertex.CullBack,
		}),
	}
}

func (p *ForwardPass) Record(cmds command.Recorder, args render.Args, scene object.T) {
	forwardMeshes := object.Query[mesh.T]().
		Where(isDrawForward).
		Collect(scene)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf)
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

func isDrawForward(m mesh.T) bool {
	return m.Mode() == mesh.Forward
}
