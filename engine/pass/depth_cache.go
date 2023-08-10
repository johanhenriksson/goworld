package pass

import (
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type DepthMatCache struct {
	app    vulkan.App
	pass   renderpass.T
	frames int
}

func NewDepthMaterialCache(app vulkan.App, pass renderpass.T, frames int) MaterialCache {
	return cache.New[*material.Def, []Material](&DepthMatCache{
		app:    app,
		pass:   pass,
		frames: frames,
	})
}

func (m *DepthMatCache) Name() string { return "DepthMaterials" }

func (m *DepthMatCache) Instantiate(def *material.Def, callback func([]Material)) {
	if def == nil {
		def = &material.Def{}
	}

	desc := &BasicDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   2000,
		},
	}

	// read vertex pointers from vertex format
	pointers := vertex.ParsePointers(def.VertexFormat)

	// fetch shader from cache
	shader := m.app.Shaders().Fetch(shader.NewRef("depth"))

	// create material
	mat := material.New(
		m.app.Device(),
		material.Args{
			Shader:     shader,
			Pass:       m.pass,
			Subpass:    MainSubpass,
			Pointers:   pointers,
			CullMode:   vertex.CullBack,
			DepthTest:  true,
			DepthWrite: true,
			DepthFunc:  core1_0.CompareOpLess,
			DepthClamp: def.DepthClamp,
			Primitive:  def.Primitive,
		},
		desc)

	instances := make([]Material, m.frames)
	for i := range instances {
		instance := mat.Instantiate(m.app.Pool())
		instances[i] = &BasicMaterial{
			id:       def.Hash(),
			Instance: instance,
			Objects:  NewObjectBuffer(desc.Objects.Size),
			Meshes:   m.app.Meshes(),
		}
	}

	callback(instances)
}

func (m *DepthMatCache) Destroy() {
}

func (m *DepthMatCache) Delete(mat []Material) {
	mat[0].Destroy()
}
