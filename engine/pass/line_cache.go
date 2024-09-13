package pass

import (
	"fmt"

	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type LineMatCache struct {
	app    engine.App
	pass   *renderpass.Renderpass
	frames int
}

func NewLineMaterialCache(app engine.App, pass *renderpass.Renderpass, frames int) MaterialCache {
	return cache.New[*material.Def, []Material](&LineMatCache{
		app:    app,
		pass:   pass,
		frames: frames,
	})
}

func (m *LineMatCache) Name() string { return "LineMaterials" }

func (m *LineMatCache) Instantiate(def *material.Def, callback func([]Material)) {
	if def == nil {
		def = material.Lines()
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
	shader := m.app.Shaders().Fetch(shader.Ref(def.Shader))

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
		instances[i] = &BasicMaterial{
			id:       def.Hash(),
			Instance: instance,
			Objects:  NewObjectBuffer(desc.Objects.Size),
			Meshes:   m.app.Meshes(),
			Commands: command.NewIndirectDrawBuffer(m.app.Device(),
				fmt.Sprintf("LineCommands:%d", i),
				desc.Objects.Size),
		}
	}

	callback(instances)
}

func (m *LineMatCache) Destroy() {
}

func (m *LineMatCache) Delete(mat []Material) {
	for _, m := range mat {
		m.Destroy()
	}
}
