package pass

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/engine/renderer/uniform"
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
	gbuffer GeometryBuffer
	target  vulkan.Target
	pass    renderpass.T
	fbuf    framebuffer.T
	fwdmat  material.Standard
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

	fwdmat := material.FromDef(
		target.Device(),
		pool,
		pass,
		&material.Def{
			Shader:       "vk/forward",
			Subpass:      "forward",
			VertexFormat: vertex.C{},
		})

	return &ForwardPass{
		gbuffer: gbuffer,
		target:  target,
		pass:    pass,

		fbuf:   fbuf,
		fwdmat: fwdmat,
	}
}

func (p *ForwardPass) Record(cmds command.Recorder, args render.Args, scene object.T) {
	camera := uniform.Camera{
		Proj:        args.Projection,
		View:        args.View,
		ViewProj:    args.VP,
		ProjInv:     args.Projection.Invert(),
		ViewInv:     args.View.Invert(),
		ViewProjInv: args.VP.Invert(),
		Eye:         args.Position,
	}

	p.fwdmat.Descriptors().Camera.Set(camera)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf)
		p.fwdmat.Bind(cmd)
	})

	forwardMeshes := query.New[mesh.T]().
		Where(isDrawForward).
		Collect(scene)

	for index, mesh := range forwardMeshes {
		p.DrawForward(cmds, index, mesh)
	}

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})
}

func (p *ForwardPass) DrawForward(cmds command.Recorder, index int, mesh mesh.T) {
	vkmesh := p.target.Meshes().Fetch(mesh.Mesh())
	if vkmesh == nil {
		return
	}

	p.fwdmat.Descriptors().Objects.Set(index, uniform.Object{
		Model: mesh.Transform().World(),
	})

	cmds.Record(func(cmd command.Buffer) {
		vkmesh.Draw(cmd, index)
	})
}

func (p *ForwardPass) Name() string {
	return "Forward"
}

func (p *ForwardPass) Destroy() {
	p.fbuf.Destroy()
	p.pass.Destroy()
	p.fwdmat.Material().Destroy()
}

func isDrawForward(m mesh.T) bool {
	return m.Mode() == mesh.Forward
}
