package pass

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/light"
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
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type ShadowmapLookupFn func(light.T, int) *texture.Texture

type Shadowpass struct {
	app    engine.App
	target engine.Target
	pass   *renderpass.Renderpass
	size   int

	layout     *pipeline.Layout
	descLayout *descriptor.Layout[*BasicDescriptors]
	objects    *uniform.ObjectBuffer
	plan       *RenderPlan
	commands   []*command.IndirectDrawBuffer

	// should be replaced with a proper cache that will evict unused maps
	shadowmaps map[light.T]Shadowmap

	meshes     cache.MeshCache
	pipelines  cache.PipelineCache
	lightQuery *object.Query[light.T]
	meshQuery  *object.Query[mesh.Mesh]
}

type Shadowmap struct {
	Cascades []Cascade
}

func (s *Shadowmap) Destroy() {
	for _, cascade := range s.Cascades {
		cascade.Destroy()
	}
}

type Cascade struct {
	Texture     *texture.Texture
	Frame       *framebuffer.Framebuffer
	Descriptors []*BasicDescriptors
}

func (c *Cascade) Destroy() {
	c.Texture.Destroy()
	c.Frame.Destroy()
	for _, desc := range c.Descriptors {
		desc.Destroy()
	}
}

func NewShadowPass(app engine.App, target engine.Target) *Shadowpass {
	pass := renderpass.New(app.Device(), renderpass.Args{
		Name: "Shadow",
		DepthAttachment: &attachment.Depth{
			Image:         attachment.NewImage("shadowmap", core1_0.FormatD32SignedFloat, core1_0.ImageUsageDepthStencilAttachment|core1_0.ImageUsageInputAttachment|core1_0.ImageUsageSampled),
			LoadOp:        core1_0.AttachmentLoadOpClear,
			StencilLoadOp: core1_0.AttachmentLoadOpClear,
			StoreOp:       core1_0.AttachmentStoreOpStore,
			FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
			ClearDepth:    1,
		},
		Subpasses: []renderpass.Subpass{
			{
				Name:  MainSubpass,
				Depth: true,
			},
		},
		Dependencies: []renderpass.SubpassDependency{
			{
				Src:   renderpass.ExternalSubpass,
				Dst:   MainSubpass,
				Flags: core1_0.DependencyByRegion,

				// External passes must finish reading depth textures in fragment shaders
				SrcStageMask:  core1_0.PipelineStageEarlyFragmentTests | core1_0.PipelineStageLateFragmentTests,
				SrcAccessMask: core1_0.AccessDepthStencilAttachmentRead,

				// Before we can write to the depth buffer
				DstStageMask:  core1_0.PipelineStageEarlyFragmentTests | core1_0.PipelineStageLateFragmentTests,
				DstAccessMask: core1_0.AccessDepthStencilAttachmentWrite,
			},
			{
				Src:   MainSubpass,
				Dst:   renderpass.ExternalSubpass,
				Flags: core1_0.DependencyByRegion,

				// The shadow pass must finish writing the depth attachment
				SrcStageMask:  core1_0.PipelineStageEarlyFragmentTests | core1_0.PipelineStageLateFragmentTests,
				SrcAccessMask: core1_0.AccessDepthStencilAttachmentWrite,

				// Before it can be used as a shadow map texture in a fragment shader
				DstStageMask:  core1_0.PipelineStageFragmentShader,
				DstAccessMask: core1_0.AccessShaderRead,
			},
		},
	})

	maxObjects := 1000
	descLayout := descriptor.NewLayout(app.Device(), "Shadows", &BasicDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   maxObjects,
		},
	})
	layout := pipeline.NewLayout(app.Device(), []descriptor.SetLayout{descLayout}, []pipeline.PushConstant{})

	objects := uniform.NewObjectBuffer(maxObjects)
	pipelines := cache.NewPipelineCache(app.Device(), app.Shaders(), pass, layout)

	commands := make([]*command.IndirectDrawBuffer, target.Frames())
	for i := range commands {
		commands[i] = command.NewIndirectDrawBuffer(app.Device(), "Shadows", objects.Size())
	}

	return &Shadowpass{
		app:        app,
		target:     target,
		pass:       pass,
		shadowmaps: make(map[light.T]Shadowmap),
		size:       2048,

		layout:     layout,
		descLayout: descLayout,
		objects:    objects,
		commands:   commands,
		plan:       NewRenderPlan(),

		pipelines:  pipelines,
		meshes:     app.Meshes(),
		meshQuery:  object.NewQuery[mesh.Mesh](),
		lightQuery: object.NewQuery[light.T](),
	}
}

func (p *Shadowpass) Name() string {
	return "Shadow"
}

func (p *Shadowpass) createShadowmap(light light.T) Shadowmap {
	log.Println("creating shadowmap for", light.Name())

	cascades := make([]Cascade, light.Shadowmaps())
	for i := range cascades {
		key := fmt.Sprintf("%s-%d", object.Key("light", light), i)
		fbuf, err := framebuffer.New(p.app.Device(), key, p.size, p.size, p.pass)
		if err != nil {
			panic(err)
		}

		// the frame buffer object will allocate a new depth image for us
		view := fbuf.Attachment(attachment.DepthName)
		tex, err := texture.FromView(p.app.Device(), key, view, texture.Args{
			Aspect: core1_0.ImageAspectDepth,
		})
		if err != nil {
			panic(err)
		}

		cascades[i].Texture = tex
		cascades[i].Frame = fbuf

		// each light cascade needs its own descriptors projection
		// todo: share object descriptors between cascades
		cascades[i].Descriptors = p.descLayout.InstantiateMany(p.app.Pool(), p.target.Frames())
	}

	shadowmap := Shadowmap{
		Cascades: cascades,
	}
	p.shadowmaps[light] = shadowmap
	return shadowmap
}

func (p *Shadowpass) fetch(mesh mesh.Mesh) (*cache.GpuMesh, *cache.Pipeline, bool) {
	gpuMesh, meshReady := p.meshes.TryFetch(mesh.Mesh())
	if !meshReady {
		return nil, nil, false
	}

	def := *mesh.Material()
	def.Shader = "pass/shadow"
	def.CullMode = vertex.CullFront
	def.DepthTest = true
	def.DepthWrite = true
	def.DepthClamp = true

	mat, matReady := p.pipelines.TryFetch(&def)
	if !matReady {
		return nil, nil, false
	}

	return gpuMesh, mat, true
}

func (p *Shadowpass) Record(cmds command.Recorder, args draw.Args, scene object.Component) {
	indirect := p.commands[args.Frame]

	lights := p.lightQuery.
		Reset().
		Where(func(lit light.T) bool { return lit.Type() == light.TypeDirectional && lit.CastShadows() }).
		Collect(scene)

	meshes := p.meshQuery.
		Reset().
		Where(castsShadows).
		Collect(scene)

	p.plan.Clear()
	p.objects.Reset()

	// record render plan with all shadow casters
	for _, meshObject := range meshes {
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

	// todo: frustum cull meshes using light frustum

	for _, light := range lights {
		shadowmap, mapExists := p.shadowmaps[light]
		if !mapExists {
			shadowmap = p.createShadowmap(light)
		}

		for index, cascade := range shadowmap.Cascades {
			camera := light.ShadowProjection(index)
			frame := cascade.Frame

			// update descriptors
			desc := cascade.Descriptors[args.Frame]
			desc.Camera.Set(camera)
			p.objects.Flush(desc.Objects)

			cmds.Record(func(cmd *command.Buffer) {
				cmd.CmdBeginRenderPass(p.pass, frame)
				cmd.CmdBindGraphicsDescriptor(p.layout, 0, desc)
				p.plan.Draw(cmd, indirect)
				cmd.CmdEndRenderPass()
			})
		}
	}
}

func castsShadows(m mesh.Mesh) bool {
	return m.CastShadows()
}

func (p *Shadowpass) Shadowmap(light light.T, cascade int) *texture.Texture {
	if shadowmap, exists := p.shadowmaps[light]; exists {
		return shadowmap.Cascades[cascade].Texture
	}
	return nil
}

func (p *Shadowpass) Destroy() {
	for _, shadowmap := range p.shadowmaps {
		shadowmap.Destroy()
	}
	p.shadowmaps = nil

	p.pass.Destroy()
	p.pass = nil

	for _, commands := range p.commands {
		commands.Destroy()
	}
	p.commands = nil

	p.layout.Destroy()
	p.layout = nil

	p.descLayout.Destroy()
	p.descLayout = nil

	p.pipelines.Destroy()
	p.pipelines = nil
}
