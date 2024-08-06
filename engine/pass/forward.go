package pass

import (
	"github.com/johanhenriksson/goworld/core/light"
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
	target vulkan.Target
	app    vulkan.App
	pass   renderpass.T
	fbuf   framebuffer.Array

	materials  MaterialCache
	meshQuery  *object.Query[mesh.Mesh]
	lightQuery *object.Query[light.T]
}

var _ Pass = &ForwardPass{}

func NewForwardPass(
	app vulkan.App,
	target vulkan.Target,
	depth vulkan.Target,
	shadows Shadow,
) *ForwardPass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Forward",
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				LoadOp:        core1_0.AttachmentLoadOpLoad,
				StoreOp:       core1_0.AttachmentStoreOpStore,
				InitialLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
				Blend:         attachment.BlendMultiply,

				Image: attachment.FromImageArray(target.Surfaces()),
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpLoad,
			StencilLoadOp: core1_0.AttachmentLoadOpLoad,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			InitialLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,

			Image: attachment.FromImageArray(depth.Surfaces()),
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{OutputAttachment},
			},
		},
	})

	fbuf, err := framebuffer.NewArray(target.Frames(), app.Device(), "forward", target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	return &ForwardPass{
		target: target,
		app:    app,
		pass:   pass,
		fbuf:   fbuf,

		materials:  NewForwardMaterialCache(app, pass, target.Frames(), shadows.Shadowmap),
		meshQuery:  object.NewQuery[mesh.Mesh](),
		lightQuery: object.NewQuery[light.T](),
	}
}

func (p *ForwardPass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	cam := CameraFromArgs(args)
	lights := p.lightQuery.Reset().Collect(scene)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Frame])
	})

	// opaque pass
	opaque := p.meshQuery.
		Reset().
		Where(isDrawForward).
		Where(isTransparent(false)).
		Collect(scene)
	groups := MaterialGroups(p.materials, args.Frame, opaque)
	groups.Draw(cmds, cam, lights)

	// transparent pass
	transparent := p.meshQuery.
		Reset().
		Where(isDrawForward).
		Where(isTransparent(true)).
		Collect(scene)
	groups = DepthSortGroups(p.materials, args.Frame, cam, transparent)
	groups.Draw(cmds, cam, lights)

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

func isTransparent(transparent bool) func(m mesh.Mesh) bool {
	return func(m mesh.Mesh) bool {
		if mat := m.Material(); mat != nil {
			return m.Material().Transparent == transparent
		}
		return false
	}
}
