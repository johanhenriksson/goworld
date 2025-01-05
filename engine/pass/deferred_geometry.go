package pass

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/shape"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
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

type DeferredDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Textures *descriptor.SamplerArray
}

type DeferredGeometryPass struct {
	target  engine.Target
	gbuffer GeometryBuffer
	app     engine.App
	pass    *renderpass.Renderpass
	fbuf    framebuffer.Array

	layout      *pipeline.Layout
	descLayout  *descriptor.Layout[*DeferredDescriptors]
	descriptors []*DeferredDescriptors
	textures    *cache.SamplerCache
	objects     *uniform.ObjectBuffer
	plan        *RenderPlan
	commands    []*command.IndirectDrawBuffer

	meshes    cache.MeshCache
	pipelines cache.PipelineCache
	meshQuery *object.Query[mesh.Mesh]
}

var _ draw.Pass = (*DeferredGeometryPass)(nil)

func NewDeferredGeometryPass(
	app engine.App,
	depth engine.Target,
	gbuffer GeometryBuffer,
) *DeferredGeometryPass {
	maxTextures := 100
	maxObjects := 1000

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

	// todo: these could probably be global descriptors
	// pass descriptor layout
	descLayout := descriptor.NewLayout(app.Device(), "DeferredGeometry", &DeferredDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   maxObjects,
		},
		Textures: &descriptor.SamplerArray{
			Stages: core1_0.StageFragment,
			Count:  maxTextures,
		},
	})
	descriptors := descLayout.InstantiateMany(app.Pool(), gbuffer.Frames())
	layout := pipeline.NewLayout(app.Device(), []descriptor.SetLayout{descLayout}, nil)

	textures := cache.NewSamplerCache(app.Textures(), maxTextures)
	objects := uniform.NewObjectBuffer(maxObjects)
	pipelines := cache.NewPipelineCache(app.Device(), app.Shaders(), pass, layout)

	commands := make([]*command.IndirectDrawBuffer, gbuffer.Frames())
	for i := range commands {
		commands[i] = command.NewIndirectDrawBuffer(app.Device(), "Deferred", objects.Size())
	}

	app.Textures().Fetch(color.White)

	return &DeferredGeometryPass{
		gbuffer: gbuffer,
		pass:    pass,
		fbuf:    fbuf,

		layout:      layout,
		descLayout:  descLayout,
		descriptors: descriptors,
		objects:     objects,
		textures:    textures,
		commands:    commands,
		plan:        NewRenderPlan(),

		pipelines: pipelines,
		meshes:    app.Meshes(),
		meshQuery: object.NewQuery[mesh.Mesh](),
	}
}

func (p *DeferredGeometryPass) fetch(mesh mesh.Mesh) (*cache.GpuMesh, *cache.Pipeline, bool) {
	gpuMesh, meshReady := p.meshes.TryFetch(mesh.Mesh())
	if !meshReady {
		return nil, nil, false
	}

	mat, matReady := p.pipelines.TryFetch(mesh.Material())
	if !matReady {
		return nil, nil, false
	}

	return gpuMesh, mat, true
}

func (p *DeferredGeometryPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	descriptors := p.descriptors[args.Frame]
	indirect := p.commands[args.Frame]
	framebuf := p.fbuf[args.Frame]

	// update camera descriptor
	cam := uniform.CameraFromArgs(args)
	descriptors.Camera.Set(cam)

	// clear object buffer
	p.objects.Reset()

	// collect all objects
	objects := p.meshQuery.
		Reset().
		Where(isDrawDeferred).
		Collect(scene)

	// clear render plan
	p.plan.Clear()

	for _, meshObject := range objects {
		mesh, pipeline, ready := p.fetch(meshObject)
		if !ready {
			continue
		}

		// this could happen inside the mesh cache!
		// basically *GpuMesh could be the entire uniform object
		// or even the entire object buffer similar to the sampler cache?
		textureIds := AssignMeshTextures(p.textures, meshObject, pipeline.Slots)

		objectId := p.objects.Store(uniform.Object{
			Model:    meshObject.Transform().Matrix(),
			Textures: textureIds,
			Vertices: mesh.Vertices.Address(),
			Indices:  mesh.Indices.Address(),
		})

		p.plan.Add(pipeline, RenderObject{
			Handle:  objectId,
			Indices: mesh.IndexCount,
		})
	}

	// flush descriptors
	p.objects.Flush(descriptors.Objects)
	p.textures.Flush(descriptors.Textures)

	cmds.Record(func(cmd *command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, framebuf)
		cmd.CmdBindGraphicsDescriptor(p.layout, 0, descriptors)
		p.plan.Draw(cmd, indirect)
		cmd.CmdEndRenderPass()
	})
}

func (p *DeferredGeometryPass) Name() string {
	return "Deferred"
}

func (p *DeferredGeometryPass) Destroy() {
	p.fbuf.Destroy()
	p.pass.Destroy()
	for _, desc := range p.descriptors {
		desc.Destroy()
	}
	for _, commands := range p.commands {
		commands.Destroy()
	}
	p.layout.Destroy()
	p.descLayout.Destroy()
	p.pipelines.Destroy()
}

func isDrawDeferred(m mesh.Mesh) bool {
	if ref := m.Mesh(); ref != nil {
		if mat := m.Material(); mat != nil {
			return mat.Pass == material.Deferred
		}
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
