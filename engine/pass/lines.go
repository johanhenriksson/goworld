package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	lineShape "github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type LinePass struct {
	app       vulkan.App
	target    vulkan.Target
	pass      renderpass.T
	fbuf      framebuffer.Array
	materials MatCache
	meshQuery *object.Query[mesh.Mesh]
}

func NewLinePass(app vulkan.App, target vulkan.Target, depth vulkan.Target) *LinePass {
	log.Println("create line pass")

	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Lines",
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				Image:         attachment.FromImageArray(target.Surfaces()),
				LoadOp:        core1_0.AttachmentLoadOpLoad,
				StoreOp:       core1_0.AttachmentStoreOpStore,
				InitialLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
				Blend:         attachment.BlendMix,
			},
		},
		DepthAttachment: &attachment.Depth{
			Image:         attachment.FromImageArray(depth.Surfaces()),
			LoadOp:        core1_0.AttachmentLoadOpLoad,
			InitialLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
			FinalLayout:   core1_0.ImageLayoutDepthStencilAttachmentOptimal,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{OutputAttachment},
			},
		},
	})

	fbufs, err := framebuffer.NewArray(target.Frames(), app.Device(), "lines", target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	lineShape.Debug.Setup(target.Frames())

	return &LinePass{
		app:       app,
		target:    target,
		pass:      pass,
		fbuf:      fbufs,
		materials: NewLineMaterialCache(app, pass, target.Frames()),
		meshQuery: object.NewQuery[mesh.Mesh](),
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

	// debug lines
	debug := lineShape.Debug.Fetch()
	lines = append(lines, debug)

	cam := CameraFromArgs(args)
	groups := MaterialGroups(p.materials, args.Context.Index, lines)
	groups.Draw(cmds, cam, nil)

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

func isDrawLines(m mesh.Mesh) bool {
	return m.Primitive() == vertex.Lines
}
