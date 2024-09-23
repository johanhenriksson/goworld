package pass

import (
	"fmt"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type ForwardMatCache struct {
	app     engine.App
	pass    *renderpass.Renderpass
	shaders cache.ShaderCache
	meshes  cache.MeshCache
	frames  int

	pipeLayout *pipeline.Layout
	descLayout *descriptor.Layout[*ForwardDescriptors]

	descriptors []*ForwardDescriptors
	objects     *ObjectBuffer
	lights      *LightBuffer
	shadows     *ShadowCache
	textures    cache.SamplerCache
}

func NewForwardMaterialCache(app engine.App, pass *renderpass.Renderpass, frames int, textures cache.SamplerCache, objects *ObjectBuffer, lights *LightBuffer, shadows *ShadowCache) MaterialCache {

	// these are actually the global/pass descriptors
	descLayout := descriptor.NewLayout(app.Device(), "Forward", &ForwardDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   objects.Size(),
		},
		Lights: &descriptor.Storage[uniform.Light]{
			Stages: core1_0.StageAll,
			Size:   lights.Size(),
		},
		Textures: &descriptor.SamplerArray{
			Stages: core1_0.StageFragment,
			Count:  textures.Size(),
		},
	})
	pipeLayout := pipeline.NewLayout(app.Device(), []descriptor.SetLayout{descLayout}, []pipeline.PushConstant{})

	descriptors := descLayout.InstantiateMany(app.Pool(), frames)

	return cache.New[*material.Def, []Material](&ForwardMatCache{
		app:        app,
		pass:       pass,
		shaders:    app.Shaders(),
		meshes:     app.Meshes(),
		frames:     frames,
		descLayout: descLayout,
		pipeLayout: pipeLayout,

		descriptors: descriptors,
		objects:     objects,
		lights:      lights,
		textures:    textures,
		shadows:     shadows,
	})
}

func (m *ForwardMatCache) Name() string { return "ForwardMaterials" }

func (m *ForwardMatCache) Instantiate(def *material.Def, callback func([]Material)) {
	if def == nil {
		def = material.StandardForward()
	}

	// read vertex pointers from vertex format
	pointers := vertex.ParsePointers(def.VertexFormat)

	// fetch shader from cache
	shader := m.shaders.Fetch(shader.Ref(def.Shader))

	// create material
	pipe := pipeline.New(
		m.app.Device(),
		pipeline.Args{
			Shader:     shader,
			Layout:     m.pipeLayout,
			Pass:       m.pass,
			Subpass:    MainSubpass,
			Pointers:   pointers,
			DepthTest:  def.DepthTest,
			DepthWrite: def.DepthWrite,
			DepthClamp: def.DepthClamp,
			DepthFunc:  def.DepthFunc,
			Primitive:  def.Primitive,
			CullMode:   def.CullMode,
		})

	instances := make([]Material, m.frames)
	for i := range instances {
		instances[i] = &ForwardMaterial{
			// material instance details:

			id:          def.Hash(),
			Pipeline:    pipe,
			Descriptors: m.descriptors[i], // per-material descriptors

			// this could come from somewhere else:

			Commands: command.NewIndirectDrawBuffer(m.app.Device(),
				fmt.Sprintf("ForwardCommands:%d", i),
				m.descriptors[i].Objects.Size),

			Objects:  m.objects,
			Lights:   m.lights,
			Shadows:  m.shadows,
			Textures: m.textures,
			Meshes:   m.meshes,
		}
	}

	callback(instances)
}

func (m *ForwardMatCache) Destroy() {
	m.textures.Destroy()
	for _, desc := range m.descriptors {
		desc.Destroy()
	}
	m.pipeLayout.Destroy()
	m.descLayout.Destroy()
}

func (m *ForwardMatCache) Delete(mat []Material) {
	for _, m := range mat {
		m.Destroy()
	}
}
