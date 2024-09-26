package pass

import (
	"sort"

	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type ForwardDescriptors struct {
	descriptor.Set
	Camera   *descriptor.Uniform[uniform.Camera]
	Objects  *descriptor.Storage[uniform.Object]
	Lights   *descriptor.Storage[uniform.Light]
	Textures *descriptor.SamplerArray
}

type ForwardPass struct {
	target engine.Target
	pass   *renderpass.Renderpass
	fbuf   framebuffer.Array

	layout      *pipeline.Layout
	descLayout  *descriptor.Layout[*ForwardDescriptors]
	descriptors []*ForwardDescriptors
	textures    cache.SamplerCache
	objects     *ObjectBuffer
	lights      *LightBuffer
	shadows     *ShadowCache
	plan        *RenderPlan
	commands    []*command.IndirectDrawBuffer

	meshes     cache.MeshCache
	pipelines  cache.PipelineCache
	meshQuery  *object.Query[mesh.Mesh]
	lightQuery *object.Query[light.T]
}

var _ draw.Pass = &ForwardPass{}

func NewForwardPass(
	app engine.App,
	target engine.Target,
	depth engine.Target,
	shadowPass *Shadowpass,
) *ForwardPass {
	// todo: arguments/settings
	maxLights := 256
	maxTextures := 100
	maxObjects := 1000

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

	// todo: these could probably be global descriptors
	// pass descriptor layout
	descLayout := descriptor.NewLayout(app.Device(), "Forward", &ForwardDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   maxObjects,
		},
		Lights: &descriptor.Storage[uniform.Light]{
			Stages: core1_0.StageAll,
			Size:   maxLights,
		},
		Textures: &descriptor.SamplerArray{
			Stages: core1_0.StageFragment,
			Count:  maxTextures,
		},
	})
	descriptors := descLayout.InstantiateMany(app.Pool(), target.Frames())
	layout := pipeline.NewLayout(app.Device(), []descriptor.SetLayout{descLayout}, nil)

	textures := cache.NewSamplerCache(app.Textures(), maxTextures)
	objects := NewObjectBuffer(maxObjects)
	lights := NewLightBuffer(maxLights)
	shadows := NewShadowCache(textures, shadowPass.Shadowmap)
	pipelines := cache.NewPipelineCache(app.Device(), app.Shaders(), pass, layout)

	commands := make([]*command.IndirectDrawBuffer, target.Frames())
	for i := range commands {
		commands[i] = command.NewIndirectDrawBuffer(app.Device(), "Forward", objects.Size())
	}

	return &ForwardPass{
		target: target,
		pass:   pass,
		fbuf:   fbuf,

		layout:      layout,
		descLayout:  descLayout,
		descriptors: descriptors,
		objects:     objects,
		lights:      lights,
		textures:    textures,
		shadows:     shadows,
		commands:    commands,
		plan:        NewRenderPlan(),

		pipelines:  pipelines,
		meshes:     app.Meshes(),
		meshQuery:  object.NewQuery[mesh.Mesh](),
		lightQuery: object.NewQuery[light.T](),
	}
}

func (p *ForwardPass) fetch(mesh mesh.Mesh) (*cache.GpuMesh, *cache.Pipeline, bool) {
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

type transparentObject struct {
	Mesh     mesh.Mesh
	GpuMesh  *cache.GpuMesh
	Pipeline *cache.Pipeline
}

func (p *ForwardPass) depthSort(objects []mesh.Mesh, eye vec3.T) []transparentObject {
	meshes := make([]transparentObject, 0, len(objects))
	for _, meshObject := range objects {
		mesh, pipeline, ready := p.fetch(meshObject)
		if !ready {
			continue
		}

		meshes = append(meshes, transparentObject{
			Mesh:     meshObject,
			GpuMesh:  mesh,
			Pipeline: pipeline,
		})
	}
	sort.SliceStable(meshes, func(i, j int) bool {
		// return true if meshes[i] is closer than meshes[j]
		pi, pj := meshes[i].Mesh.Transform().WorldPosition(), meshes[j].Mesh.Transform().WorldPosition()
		bi, bj := meshes[i].GpuMesh.Bounds(), meshes[j].GpuMesh.Bounds()
		di := vec3.Distance(eye, bi.Center.Add(pi)) - bi.Radius
		dj := vec3.Distance(eye, bj.Center.Add(pj)) - bj.Radius
		return di < dj
	})
	return meshes
}

func (p *ForwardPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	descriptors := p.descriptors[args.Frame]
	indirect := p.commands[args.Frame]
	framebuf := p.fbuf[args.Frame]

	// update camera descriptor
	cam := uniform.CameraFromArgs(args)
	descriptors.Camera.Set(cam)

	// fill light buffer
	lights := p.lightQuery.Reset().Collect(scene)
	p.lights.Reset()
	for _, lit := range lights {
		p.lights.Store(lit.LightData(p.shadows))
	}

	// clear object buffer
	p.objects.Reset()

	// opaque pass
	opaqueQuery := p.meshQuery.
		Reset().
		Where(isDrawForward).
		Where(isTransparent(false)).
		Collect(scene)

	// clear render plan
	p.plan.Clear()

	for _, meshObject := range opaqueQuery {
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

	// transparent pass
	transparentQuery := p.meshQuery.
		Reset().
		Where(isDrawForward).
		Where(isTransparent(true)).
		Collect(scene)

	// depth sort transparent meshes
	transparentObjects := p.depthSort(transparentQuery, args.Camera.Position)

	for _, t := range transparentObjects {
		textureIds := AssignMeshTextures(p.textures, t.Mesh, t.Pipeline.Slots)

		objectId := p.objects.Store(uniform.Object{
			Model:    t.Mesh.Transform().Matrix(),
			Textures: textureIds,
			Vertices: t.GpuMesh.Vertices.Address(),
			Indices:  t.GpuMesh.Indices.Address(),
		})

		p.plan.AddOrdered(t.Pipeline, RenderObject{
			Handle:  objectId,
			Indices: t.GpuMesh.IndexCount,
		})
	}

	// flush descriptors
	p.lights.Flush(descriptors.Lights)
	p.objects.Flush(descriptors.Objects)
	p.textures.Flush(descriptors.Textures)

	//
	// phase 2: record commands
	//

	cmds.Record(func(cmd *command.Buffer) {
		indirect.Reset()
		cmd.CmdBeginRenderPass(p.pass, framebuf)
		cmd.CmdBindGraphicsDescriptor(p.layout, 0, descriptors)
		p.plan.Draw(cmd, indirect)
		cmd.CmdEndRenderPass()
	})
}

func (p *ForwardPass) Name() string {
	return "Forward"
}

func (p *ForwardPass) Destroy() {
	p.textures.Destroy()
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

func isDrawForward(m mesh.Mesh) bool {
	if ref := m.Mesh(); ref != nil {
		if mat := m.Material(); mat != nil {
			return mat.Pass == material.Forward
		}
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
