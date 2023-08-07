package pass

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type DepthPass struct {
	gbuffer GeometryBuffer
	app     vulkan.App
	pass    renderpass.T
	fbuf    framebuffer.Array

	materials *MeshSorter[*DepthMatData]
	meshQuery *object.Query[mesh.Mesh]
}

var _ Pass = &ForwardPass{}

func NewDepthPass(
	app vulkan.App,
	depth vulkan.Target,
	gbuffer GeometryBuffer,
) *DepthPass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Depth",
		ColorAttachments: []attachment.Color{
			{
				Name:        NormalsAttachment,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,

				Image: attachment.FromImageArray(gbuffer.Normal()),
			},
			{
				Name:        PositionAttachment,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,

				Image: attachment.FromImageArray(gbuffer.Position()),
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpClear,
			StencilLoadOp: core1_0.AttachmentLoadOpClear,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			ClearDepth:    1,

			Image: attachment.FromImageArray(depth.Surfaces()),
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{NormalsAttachment, PositionAttachment},
			},
		},
	})

	fbuf, err := framebuffer.NewArray(gbuffer.Frames(), app.Device(), "depth", gbuffer.Width(), gbuffer.Height(), pass)
	if err != nil {
		panic(err)
	}

	mats := NewMeshSorter(app, gbuffer.Frames(), NewDepthMaterialMaker(app, pass))

	return &DepthPass{
		gbuffer: gbuffer,
		app:     app,
		pass:    pass,
		fbuf:    fbuf,

		materials: mats,
		meshQuery: object.NewQuery[mesh.Mesh](),
	}
}

func (p *DepthPass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	forwardMeshes := p.meshQuery.
		Reset().
		Where(isDrawForward). // todo: dont include transparent objects
		Collect(scene)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Context.Index])
	})

	p.materials.Draw(cmds, args, forwardMeshes, nil)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})
}

func (p *DepthPass) Name() string {
	return "Depth"
}

func (p *DepthPass) Destroy() {
	p.fbuf.Destroy()
	p.pass.Destroy()
	p.materials.Destroy()
}
