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
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type ForwardPass struct {
	target  vulkan.Target
	gbuffer GeometryBuffer
	app     vulkan.App
	pass    renderpass.T
	fbuf    framebuffer.Array

	materials *MaterialSorter
	meshQuery *object.Query[mesh.Mesh]
}

var _ Pass = &ForwardPass{}

func NewForwardPass(
	app vulkan.App,
	target vulkan.Target,
	depth vulkan.Target,
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

				Image: attachment.FromImageArray(target.Surfaces()),
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

			Image: attachment.FromImageArray(depth.Surfaces()),
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{OutputAttachment, NormalsAttachment, PositionAttachment},
			},
		},
	})

	fbuf, err := framebuffer.NewArray(target.Frames(), app.Device(), "forward", target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	return &ForwardPass{
		target:  target,
		gbuffer: gbuffer,
		app:     app,
		pass:    pass,
		fbuf:    fbuf,

		materials: NewMaterialSorter(app, target, pass, material.ColoredForward()),
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
	if mat := m.Material(); mat != nil {
		return mat.Pass == material.Forward
	}
	return false
}
