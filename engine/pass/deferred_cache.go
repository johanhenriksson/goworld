package pass

import (
	"fmt"

	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type DeferredMatCache struct {
	app    vulkan.App
	pass   renderpass.T
	frames int
}

func NewDeferredMaterialCache(app vulkan.App, pass renderpass.T, frames int) MaterialCache {
	return cache.New[*material.Def, []Material](&DeferredMatCache{
		app:    app,
		pass:   pass,
		frames: frames,
	})
}

func (m *DeferredMatCache) Name() string { return "DeferredMaterials" }

func (m *DeferredMatCache) Instantiate(def *material.Def, callback func([]Material)) {
	if def == nil {
		def = material.StandardDeferred()
	}

	desc := &DeferredDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   2000,
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
		instances[i] = &DeferredMaterial{
			id:       def.Hash(),
			Instance: instance,
			Objects:  NewObjectBuffer(instance.Descriptors().Objects.Size),
			Textures: textures,
			Meshes:   m.app.Meshes(),
			Commands: command.NewIndirectDrawBuffer(m.app.Device(),
				fmt.Sprintf("DeferredCommands:%d", i),
				instance.Descriptors().Objects.Size),
		}
	}

	callback(instances)
}

func (m *DeferredMatCache) Destroy() {

}

func (m *DeferredMatCache) Delete(mat []Material) {
	for _, m := range mat {
		m.Destroy()
	}
}
