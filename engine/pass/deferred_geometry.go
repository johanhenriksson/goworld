package pass

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/shape"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

const (
	DiffuseAttachment  attachment.Name = "diffuse"
	NormalsAttachment  attachment.Name = "normals"
	PositionAttachment attachment.Name = "position"
	OutputAttachment   attachment.Name = "output"
)

type DeferredGeometryPass struct {
	target  vulkan.Target
	gbuffer GeometryBuffer
	app     vulkan.App
	pass    renderpass.T
	fbuf    framebuffer.Array

	materials *MeshSorter[*DeferredMatData]
	meshQuery *object.Query[mesh.Mesh]
}

func NewDeferredGeometryPass(
	app vulkan.App,
	depth vulkan.Target,
	gbuffer GeometryBuffer,
) *DeferredGeometryPass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Deferred Geometry",
		ColorAttachments: []attachment.Color{
			{
				Name:        DiffuseAttachment,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Image:       attachment.FromImageArray(gbuffer.Diffuse()),
			},
			{
				Name:        NormalsAttachment,
				LoadOp:      core1_0.AttachmentLoadOpLoad,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Image:       attachment.FromImageArray(gbuffer.Normal()),
			},
			{
				Name:        PositionAttachment,
				LoadOp:      core1_0.AttachmentLoadOpLoad,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Image:       attachment.FromImageArray(gbuffer.Position()),
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpLoad,
			StencilLoadOp: core1_0.AttachmentLoadOpLoad,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			Image:         attachment.FromImageArray(depth.Surfaces()),
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
				Depth: true,

				ColorAttachments: []attachment.Name{DiffuseAttachment, NormalsAttachment, PositionAttachment},
			},
		},
	})

	fbuf, err := framebuffer.NewArray(gbuffer.Frames(), app.Device(), "deferred-geometry", gbuffer.Width(), gbuffer.Height(), pass)
	if err != nil {
		panic(err)
	}

	app.Textures().Fetch(color.White)

	return &DeferredGeometryPass{
		gbuffer: gbuffer,
		app:     app,
		pass:    pass,

		fbuf: fbuf,

		materials: NewMeshSorter(app, gbuffer.Frames(), NewDeferredMaterialMaker(app, pass)),
		meshQuery: object.NewQuery[mesh.Mesh](),
	}
}

func (p *DeferredGeometryPass) Record(cmds command.Recorder, args render.Args, scene object.Component) {
	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Context.Index])
	})

	frustum := shape.FrustumFromMatrix(args.VP)

	objects := p.meshQuery.
		Reset().
		Where(isDrawDeferred).
		Where(frustumCulled(&frustum)).
		Collect(scene)

	cam := CameraFromArgs(args)
	p.materials.Draw(cmds, args.Context.Index, cam, objects, nil)

	cmds.Record(func(cmd command.Buffer) {
		cmd.CmdEndRenderPass()
	})
}

func (p *DeferredGeometryPass) Name() string {
	return "Deferred"
}

func (p *DeferredGeometryPass) Destroy() {
	p.materials.Destroy()
	p.materials = nil
	p.fbuf.Destroy()
	p.fbuf = nil
	p.pass.Destroy()
	p.pass = nil
}

func isDrawDeferred(m mesh.Mesh) bool {
	if mat := m.Material(); mat != nil {
		return mat.Pass == material.Deferred
	}
	return false
}

func frustumCulled(frustum *shape.Frustum) func(mesh.Mesh) bool {
	return func(m mesh.Mesh) bool {
		bounds := m.BoundingSphere()
		return frustum.IntersectsSphere(&bounds)
	}
}
