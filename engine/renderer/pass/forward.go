package pass

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type ForwardPass struct {
	gbuffer   GeometryBuffer
	target    vulkan.Target
	pass      renderpass.T
	pool      descriptor.Pool
	fbuf      framebuffer.T
	materials *MaterialSorter
}

func NewForwardPass(
	target vulkan.Target,
	pool descriptor.Pool,
	gbuffer GeometryBuffer,
) *ForwardPass {
	pass := renderpass.New(target.Device(), renderpass.Args{
		ColorAttachments: []attachment.Color{
			{
				Name:        OutputAttachment,
				Format:      gbuffer.Output().Format(),
				LoadOp:      vk.AttachmentLoadOpLoad,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageSampledBit,
				Blend:       attachment.BlendMix,

				Allocator: attachment.FromImageArray([]image.T{
					gbuffer.Output().Image(),
				}),
			},
			{
				Name:        NormalsAttachment,
				Format:      gbuffer.Normal().Format(),
				LoadOp:      vk.AttachmentLoadOpLoad,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit | vk.ImageUsageTransferSrcBit,

				Allocator: attachment.FromImageArray([]image.T{
					gbuffer.Normal().Image(),
				}),
			},
			{
				Name:        PositionAttachment,
				Format:      gbuffer.Position().Format(),
				LoadOp:      vk.AttachmentLoadOpLoad,
				StoreOp:     vk.AttachmentStoreOpStore,
				FinalLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				Usage:       vk.ImageUsageInputAttachmentBit | vk.ImageUsageTransferSrcBit,

				Allocator: attachment.FromImageArray([]image.T{
					gbuffer.Position().Image(),
				}),
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:        vk.AttachmentLoadOpLoad,
			StencilLoadOp: vk.AttachmentLoadOpLoad,
			StoreOp:       vk.AttachmentStoreOpStore,
			FinalLayout:   vk.ImageLayoutShaderReadOnlyOptimal,
			Usage:         vk.ImageUsageInputAttachmentBit,

			Allocator: attachment.FromImageArray([]image.T{
				gbuffer.Depth().Image(),
			}),
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  "forward",
				Depth: true,

				ColorAttachments: []attachment.Name{OutputAttachment, NormalsAttachment, PositionAttachment},
			},
		},
	})

	fbuf, err := framebuffer.New(target.Device(), target.Width(), target.Height(), pass)
	if err != nil {
		panic(err)
	}

	return &ForwardPass{
		gbuffer: gbuffer,
		target:  target,
		pass:    pass,
		pool:    pool,
		fbuf:    fbuf,

		materials: NewMaterialSorter(target, pool, pass, &material.Def{
			Shader:       "vk/forward",
			Subpass:      "forward",
			VertexFormat: vertex.C{},
			DepthTest:    true,
			DepthWrite:   true,
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
