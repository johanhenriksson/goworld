package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type LinePass struct {
	target    vulkan.Target
	pass      renderpass.T
	fbufs     framebuffer.Array
	materials *MaterialSorter
}

func NewLinePass(target vulkan.Target, geometry GeometryBuffer) *LinePass {
	log.Println("create line pass")

	depth := make([]image.T, target.Frames())
	for i := range depth {
		depth[i] = geometry.Depth().Image()
	}

	pass := renderpass.New(target.Device(), renderpass.Args{
		ColorAttachments: []attachment.Color{
			{
				Name:          OutputAttachment,
				Allocator:     attachment.FromImageArray(target.Surfaces()),
				Format:        target.SurfaceFormat(),
				LoadOp:        vk.AttachmentLoadOpLoad,
				StoreOp:       vk.AttachmentStoreOpStore,
				InitialLayout: vk.ImageLayoutPresentSrc,
				FinalLayout:   vk.ImageLayoutPresentSrc,
			},
		},
		DepthAttachment: &attachment.Depth{
			Allocator:     attachment.FromImageArray(depth),
			LoadOp:        vk.AttachmentLoadOpLoad,
			InitialLayout: vk.ImageLayoutShaderReadOnlyOptimal,
			FinalLayout:   vk.ImageLayoutDepthStencilAttachmentOptimal,
			Usage:         vk.ImageUsageInputAttachmentBit,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  OutputSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{OutputAttachment},
			},
		},
	})

	fbufs, err := framebuffer.NewArray(target.Frames(), target.Device(), target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	return &LinePass{
		target: target,
		pass:   pass,
		fbufs:  fbufs,
		materials: NewMaterialSorter(target, pass,
			&material.Def{
				Shader:       "vk/lines",
				Subpass:      OutputSubpass,
				VertexFormat: vertex.C{},
				Primitive:    vertex.Lines,
				DepthTest:    true,
			}),
	}
}

func (p *LinePass) Record(cmds command.Recorder, args render.Args, scene object.T) {
	ctx := args.Context

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbufs[ctx.Index])
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
	p.fbufs.Destroy()
	p.pass.Destroy()
	p.materials.Destroy()
}

func isDrawLines(m mesh.T) bool {
	return m.Mode() == mesh.Lines
}
