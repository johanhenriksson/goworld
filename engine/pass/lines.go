package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/uniform"
	lineShape "github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type LinePass struct {
	app       engine.App
	target    engine.Target
	pass      renderpass.T
	fbuf      framebuffer.Array
	materials MaterialCache
	meshQuery *object.Query[mesh.Mesh]
}

func NewLinePass(app engine.App, target engine.Target, depth engine.Target) *LinePass {
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

func (p *LinePass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Frame])
	})

	lines := p.meshQuery.
		Reset().
		Where(isDrawLines).
		Collect(scene)

	// debug lines
	debug := lineShape.Debug.Fetch()
	lines = append(lines, debug)

	cam := uniform.CameraFromArgs(args)
	groups := MaterialGroups(p.materials, args.Frame, lines)
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
