package pass

import (
	"log"

	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	lineShape "github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type LinePass struct {
	app    engine.App
	target engine.Target
	pass   *renderpass.Renderpass
	fbuf   framebuffer.Array

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
		commands[i] = command.NewIndirectDrawBuffer(app.Device(), "Lines", objects.Size())
	}

	lineShape.Debug.Setup(target.Frames())

	return &LinePass{
		app:    app,
		target: target,
		pass:   pass,
		fbuf:   fbufs,

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

func (p *LinePass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	descriptors := p.descriptors[args.Frame]
	indirect := p.commands[args.Frame]
	framebuf := p.fbuf[args.Frame]

	pipeline, pipeReady := p.pipelines.TryFetch(material.Lines())
	if !pipeReady {
		return
	}

	cam := uniform.CameraFromArgs(args)
	descriptors.Camera.Set(cam)

	p.objects.Reset()
	p.plan.Clear()

	lines := p.meshQuery.
		Reset().
		Where(isDrawLines).
		Collect(scene)

	// debug lines
	// todo: this causes a crash. figure out why
	// debug := lineShape.Debug.Fetch()
	// lines = append(lines, debug)

	for _, meshObject := range lines {
		mesh, meshReady := p.meshes.TryFetch(meshObject.Mesh())
		if !meshReady {
			continue
		}
		if mesh.IndexCount == 0 {
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

func (p *LinePass) Name() string {
	return "Lines"
}

func (p *LinePass) Destroy() {
	for _, desc := range p.descriptors {
		desc.Destroy()
	}
	for _, cmd := range p.commands {
		cmd.Destroy()
	}
	p.fbuf.Destroy()
	p.pass.Destroy()
	p.layout.Destroy()
	p.descLayout.Destroy()
	p.pipelines.Destroy()
}

func isDrawLines(m mesh.Mesh) bool {
	if ref := m.Mesh(); ref != nil {
		if mat := m.Material(); mat != nil {
			return m.Material().Primitive == vertex.Lines
		}
	}
	return false
}
