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

type ForwardMatCache struct {
	app    vulkan.App
	pass   renderpass.T
	lookup ShadowmapLookupFn
	frames int
}

func NewForwardMaterialCache(app vulkan.App, pass renderpass.T, frames int, lookup ShadowmapLookupFn) MaterialCache {
	return cache.New[*material.Def, []Material](&ForwardMatCache{
		app:    app,
		pass:   pass,
		lookup: lookup,
		frames: frames,
	})
}

func (m *ForwardMatCache) Name() string { return "ForwardMaterials" }

func (m *ForwardMatCache) Instantiate(def *material.Def, callback func([]Material)) {
	if def == nil {
		def = material.StandardForward()
	}

	desc := &ForwardDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   20000,
		},
		Lights: &descriptor.Storage[uniform.Light]{
			Stages: core1_0.StageAll,
			Size:   256,
		},
		Textures: &descriptor.SamplerArray{
			Stages: core1_0.StageFragment,
			Count:  100,
		},
	}

	// read vertex pointers from vertex format
	pointers := vertex.ParsePointers(def.VertexFormat)

	// fetch shader from cache
	shader := m.app.Shaders().Fetch(shader.NewRef(def.Shader))

	// create material
	mat := material.New(
		m.app.Device(),
		material.Args{
			Shader:     shader,
			Pass:       m.pass,
			Subpass:    MainSubpass,
			Pointers:   pointers,
			DepthTest:  def.DepthTest,
			DepthWrite: def.DepthWrite,
			DepthClamp: def.DepthClamp,
			DepthFunc:  def.DepthFunc,
			Primitive:  def.Primitive,
			CullMode:   def.CullMode,
		},
		desc)

	instances := make([]Material, m.frames)
	for i := range instances {
		instance := mat.Instantiate(m.app.Pool())
		textures := cache.NewSamplerCache(m.app.Textures(), instance.Descriptors().Textures)

		instances[i] = &ForwardMaterial{
			id:       def.Hash(),
			Instance: instance,
			Objects:  NewObjectBuffer(desc.Objects.Size),
			Lights:   NewLightBuffer(desc.Lights.Size),
			Shadows:  NewShadowCache(textures, m.lookup),
			Textures: textures,
			Meshes:   m.app.Meshes(),
		}
	}

	callback(instances)
}

func (m *ForwardMatCache) Destroy() {

}

func (m *ForwardMatCache) Delete(mat []Material) {
	mat[0].Destroy()
}
