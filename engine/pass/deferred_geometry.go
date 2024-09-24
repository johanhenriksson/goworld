package pass

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/shape"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"

	"github.com/vkngwrapper/core/v2/core1_0"
)

const (
	DiffuseAttachment  attachment.Name = "diffuse"
	NormalsAttachment  attachment.Name = "normals"
	PositionAttachment attachment.Name = "position"
	OutputAttachment   attachment.Name = "output"
)

type DeferredGeometryPass struct {
	target  engine.Target
	gbuffer GeometryBuffer
	app     engine.App
	pass    *renderpass.Renderpass
	fbuf    framebuffer.Array

	materials MaterialCache
	meshQuery *object.Query[mesh.Mesh]
}

var _ draw.Pass = (*DeferredGeometryPass)(nil)

func NewDeferredGeometryPass(
	app engine.App,
	depth engine.Target,
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
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Image:       attachment.FromImageArray(gbuffer.Normal()),
			},
			{
				Name:        PositionAttachment,
				LoadOp:      core1_0.AttachmentLoadOpClear,
				StoreOp:     core1_0.AttachmentStoreOpStore,
				FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
				Image:       attachment.FromImageArray(gbuffer.Position()),
			},
		},
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpLoad,
			StencilLoadOp: core1_0.AttachmentLoadOpLoad,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			InitialLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
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

		materials: NewDeferredMaterialCache(app, pass, gbuffer.Frames()),
		meshQuery: object.NewQuery[mesh.Mesh](),
	}
}

func (p *DeferredGeometryPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	cmds.Record(func(cmd *command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, p.fbuf[args.Frame])
	})

	frustum := shape.FrustumFromMatrix(args.Camera.ViewProj)

	// we would like to collect all objects at an earlier stage.
	// only traverse the scene graph once, and filter the objects based on the pass when generating indirect draw commands.
	// ideally all passes would share the same objects descriptor
	// and only changed objects would have their uniform data updated
	//
	// when traversing the scene graph, lookup meshes from the gpu cache.
	// if the mesh is not ready, skip it.
	//
	// the object buffer should contain all data required to render the object.
	// model matrix, texture ids, bounding boxes, etc.
	// this will be useful later as we move more functionality to the gpu
	//
	// problems to solve:
	// - how to have global uniform data (camera, lights, etc) that is shared between all passes?
	//   its currently deeply embedded into the material structs, and instantiated for each material
	// - once the object buffer is filled, the engine is ready to run updates for the next frame.
	//   how can we run this concurrently?

	objects := p.meshQuery.
		Reset().
		Where(isDrawDeferred).
		Where(frustumCulled(&frustum)).
		Collect(scene)

	// meshes := make([]Drawable, 0, len(objects))
	// for _, obj := range objects {
	// 	mesh, ok := p.app.Meshes().TryFetch(obj.Mesh())
	// 	if !ok {
	// 		continue
	// 	}
	//
	// 	obj.Material().TextureSlots
	// 	textureIds := AssignMeshTextures(m.Textures, obj, textures)
	//
	// 	drawable := DrawableMesh{
	// 		GpuMesh:  mesh,
	// 		model:    obj.Transform().Matrix(),
	// 		textures: textureIds,
	// 	}
	//
	// 	// frustum culling
	// 	bounds := drawable.Bounds()
	// 	if !frustum.IntersectsSphere(&bounds) {
	// 		continue
	// 	}
	//
	// 	meshes = append(meshes, drawable)
	// }

	cam := uniform.CameraFromArgs(args)
	groups := MaterialGroups(p.materials, args.Frame, objects)
	groups.Draw(cmds, cam)

	cmds.Record(func(cmd *command.Buffer) {
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
		return true
		// bounds := m.BoundingSphere()
		// return frustum.IntersectsSphere(&bounds)
	}
}
