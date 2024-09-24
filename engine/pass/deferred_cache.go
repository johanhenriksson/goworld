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

type DeferredMatCache struct {
	app        engine.App
	pass       *renderpass.Renderpass
	frames     int
	pipeLayout *pipeline.Layout
	descLayout *descriptor.Layout[*DeferredDescriptors]

	// per-frame data

	descriptors []*DeferredDescriptors
	textures    cache.SamplerCache
	objects     *ObjectBuffer
}

func NewDeferredMaterialCache(app engine.App, pass *renderpass.Renderpass, frames int) MaterialCache {
	maxTextures := 100
	maxObjects := 2000
	layout := descriptor.NewLayout(app.Device(), "Deferred", &DeferredDescriptors{
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

	descriptors := layout.InstantiateMany(app.Pool(), frames)
	textures := cache.NewSamplerCache(app.Textures(), maxTextures)
	objects := NewObjectBuffer(maxObjects)

	pipeLayout := pipeline.NewLayout(app.Device(), []descriptor.SetLayout{layout}, []pipeline.PushConstant{})

	return cache.New[*material.Def, []Material](&DeferredMatCache{
		app:        app,
		pass:       pass,
		frames:     frames,
		pipeLayout: pipeLayout,
		descLayout: layout,

		descriptors: descriptors,
		textures:    textures,
		objects:     objects,
	})
}

func (m *DeferredMatCache) Name() string { return "DeferredMaterials" }

func (m *DeferredMatCache) Instantiate(def *material.Def, callback func([]Material)) {
	if def == nil {
		def = material.StandardDeferred()
	}

	// read vertex pointers from vertex format
	pointers := vertex.ParsePointers(def.VertexFormat)

	// fetch shader from cache
	shader := m.app.Shaders().Fetch(shader.Ref(def.Shader))

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

	// all deferred materials can use the same descriptors, object buffer and sampler cache.
	// all of the materials should use the same pipeline layout, and the descriptors should be bound to it, once.
	// finally, record an indirect draw call for each material instance.

	instances := make([]Material, m.frames)
	for i := range instances {
		instances[i] = &DeferredMaterial{
			id:    def.Hash(),
			slots: shader.Textures(),

			Pipeline:    pipe,
			Descriptors: m.descriptors[i], // shared
			Objects:     m.objects,        // shared
			textures:    m.textures,       // shared
			Meshes:      m.app.Meshes(),   // maybe accessed in some other way? perhaps passed to Draw()?
			Commands: command.NewIndirectIndexedDrawBuffer(m.app.Device(),
				fmt.Sprintf("DeferredCommands:%d", i),
				m.descriptors[i].Objects.Size),
		}
	}

	callback(instances)
}

func (m *DeferredMatCache) Destroy() {
	m.textures.Destroy()
	for _, desc := range m.descriptors {
		desc.Destroy()
	}
	m.pipeLayout.Destroy()
	m.descLayout.Destroy()
}

func (m *DeferredMatCache) Delete(mat []Material) {
	for _, m := range mat {
		m.Destroy()
	}
}
