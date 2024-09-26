package pass

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type DepthPass struct {
	app   engine.App
	depth engine.Target
	pass  *renderpass.Renderpass
	fbuf  framebuffer.Array

	layout      *pipeline.Layout
	descLayout  *descriptor.Layout[*BasicDescriptors]
	descriptors []*BasicDescriptors
	objects     *ObjectBuffer
	plan        *RenderPlan
	commands    []*command.IndirectDrawBuffer

	meshes    cache.MeshCache
	pipelines cache.PipelineCache
	meshQuery *object.Query[mesh.Mesh]
}

var _ draw.Pass = &ForwardPass{}

func NewDepthPass(
	app engine.App,
	depth engine.Target,
) *DepthPass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Depth",
		DepthAttachment: &attachment.Depth{
			LoadOp:        core1_0.AttachmentLoadOpClear,
			StencilLoadOp: core1_0.AttachmentLoadOpClear,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			ClearDepth:    1,

			Image: attachment.FromImageArray(depth.Surfaces()),
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
				Depth: true,
			},
		},
	})

	fbuf, err := framebuffer.NewArray(depth.Frames(), app.Device(), "depth", depth.Width(), depth.Height(), pass)
	if err != nil {
		panic(err)
	}

	maxObjects := 1000
	descLayout := descriptor.NewLayout(app.Device(), "Depth", &BasicDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   maxObjects,
		},
	})
	descriptors := descLayout.InstantiateMany(app.Pool(), depth.Frames())
	layout := pipeline.NewLayout(app.Device(), []descriptor.SetLayout{descLayout}, []pipeline.PushConstant{})

	objects := NewObjectBuffer(maxObjects)
	pipelines := cache.NewPipelineCache(app.Device(), app.Shaders(), pass, layout)

	commands := make([]*command.IndirectDrawBuffer, depth.Frames())
	for i := range commands {
		commands[i] = command.NewIndirectDrawBuffer(app.Device(), "Depth", objects.Size())
	}

	return &DepthPass{
		app:   app,
		depth: depth,
		pass:  pass,
		fbuf:  fbuf,

		layout:      layout,
		descriptors: descriptors,
		descLayout:  descLayout,
		objects:     objects,
		commands:    commands,
		plan:        NewRenderPlan(),

		pipelines: pipelines,
		meshes:    app.Meshes(),
		meshQuery: object.NewQuery[mesh.Mesh](),
	}
}

func (p *DepthPass) fetch(mesh mesh.Mesh) (*cache.GpuMesh, *cache.Pipeline, bool) {
	def := *mesh.Material()
	def.Shader = "pass/depth"
	def.CullMode = vertex.CullBack
	def.DepthTest = true
	def.DepthWrite = true
	def.DepthClamp = true

	// only consider standard vertex format for the depth pass.
	// supporting multiple formats would require different pipelines
	if def.VertexFormat != (vertex.Vertex{}) {
		return nil, nil, false
	}

	mat, matReady := p.pipelines.TryFetch(&def)
	if !matReady {
		return nil, nil, false
	}

	gpuMesh, meshReady := p.meshes.TryFetch(mesh.Mesh())
	if !meshReady {
		return nil, nil, false
	}

	return gpuMesh, mat, true
}

func (p *DepthPass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	descriptors := p.descriptors[args.Frame]
	indirect := p.commands[args.Frame]
	framebuf := p.fbuf[args.Frame]

	cam := uniform.CameraFromArgs(args)
	descriptors.Camera.Set(cam)

	p.objects.Reset()
	p.plan.Clear()

	// todo: better strategy for picking occluders
	occluders := p.meshQuery.
		Reset().
		Where(isDrawDeferred).
		Collect(scene)

	// record render plan with all shadow casters
	for _, meshObject := range occluders {
		mesh, pipeline, ready := p.fetch(meshObject)
		if !ready {
			continue
		}

		objectId := p.objects.Store(uniform.Object{
			Model:    meshObject.Transform().Matrix(),
			Vertices: mesh.Vertices.Address(),
			Indices:  mesh.Indices.Address(),
		})

		p.plan.Add(pipeline, RenderObject{
			Handle:  objectId,
			Indices: mesh.IndexCount,
		})
	}

	p.objects.Flush(descriptors.Objects)

	cmds.Record(func(cmd *command.Buffer) {
		cmd.CmdBeginRenderPass(p.pass, framebuf)
		cmd.CmdBindGraphicsDescriptor(p.layout, 0, descriptors)
		p.plan.Draw(cmd, indirect)
		cmd.CmdEndRenderPass()
	})
}

func (p *DepthPass) Name() string {
	return "Depth"
}

func (p *DepthPass) Destroy() {
	for _, command := range p.commands {
		command.Destroy()
	}
	for _, desc := range p.descriptors {
		desc.Destroy()
	}
	p.fbuf.Destroy()
	p.pass.Destroy()
	p.layout.Destroy()
	p.descLayout.Destroy()
	p.pipelines.Destroy()
}
