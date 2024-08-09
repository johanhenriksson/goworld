package pass

import (
	"fmt"

	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/render/vulkan"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type ShadowMatCache struct {
	app    vulkan.App
	pass   renderpass.T
	frames int
}

func NewShadowMaterialMaker(app vulkan.App, pass renderpass.T, frames int) MaterialCache {
	return cache.New[*material.Def, []Material](&ShadowMatCache{
		app:    app,
		pass:   pass,
		frames: frames,
	})
}

func (m *ShadowMatCache) Name() string { return "ShadowMaterials" }

func (m *ShadowMatCache) Instantiate(def *material.Def, callback func([]Material)) {
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
	shader := m.app.Shaders().Fetch(shader.NewRef("shadow"))

	// create material
	mat := material.New(
		m.app.Device(),
		material.Args{
			Shader:     shader,
			Pass:       m.pass,
			Subpass:    MainSubpass,
			Pointers:   pointers,
			CullMode:   vertex.CullFront,
			DepthTest:  true,
			DepthWrite: true,
			DepthFunc:  core1_0.CompareOpLess,
			DepthClamp: true,
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
			Commands: command.NewIndirectDrawBuffer(m.app.Device(),
				fmt.Sprintf("ShadowCommands:%d", i),
				desc.Objects.Size),
		}
	}

	callback(instances)
}

func (m *ShadowMatCache) Destroy() {
}

func (m *ShadowMatCache) Delete(mat []Material) {
	for _, m := range mat {
		m.Destroy()
	}
}
