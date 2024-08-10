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
	app   vulkan.App
	depth vulkan.Target
	pass  renderpass.T
	fbuf  framebuffer.Array

	materials MaterialCache
	meshQuery *object.Query[mesh.Mesh]
}

var _ Pass = &ForwardPass{}

func NewDepthPass(
	app vulkan.App,
	depth vulkan.Target,
) *DepthPass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Depth",
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
			},
		},
	})

	fbuf, err := framebuffer.NewArray(depth.Frames(), app.Device(), "depth", depth.Width(), depth.Height(), pass)
	if err != nil {
		panic(err)
	}

	return &DepthPass{
		app:   app,
		depth: depth,
		pass:  pass,
		fbuf:  fbuf,

		materials: NewDepthMaterialCache(app, pass, depth.Frames()),
		meshQuery: object.NewQuery[mesh.Mesh](),
	}
}

func (p *DepthPass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	opaque := p.meshQuery.
		Reset().
		Where(isDrawDeferred).
		Collect(scene)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Frame])
	})

	cam := CameraFromArgs(args)
	groups := MaterialGroups(p.materials, args.Frame, opaque)
	groups.Draw(cmds, cam, nil)

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
