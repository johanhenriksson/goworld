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

type DepthMatCache struct {
	app    engine.App
	pass   *renderpass.Renderpass
	frames int

	descLayout *descriptor.Layout[*BasicDescriptors]
	pipeLayout *pipeline.Layout
}

func NewDepthMaterialCache(app engine.App, pass *renderpass.Renderpass, frames int) MaterialCache {
	descLayout := descriptor.NewLayout(app.Device(), "Depth", &BasicDescriptors{
		Camera: &descriptor.Uniform[uniform.Camera]{
			Stages: core1_0.StageAll,
		},
		Objects: &descriptor.Storage[uniform.Object]{
			Stages: core1_0.StageAll,
			Size:   10000,
		},
	})
	pipeLayout := pipeline.NewLayout(app.Device(), []descriptor.SetLayout{descLayout}, []pipeline.PushConstant{})
	return cache.New[*material.Def, []Material](&DepthMatCache{
		app:        app,
		pass:       pass,
		frames:     frames,
		descLayout: descLayout,
		pipeLayout: pipeLayout,
	})
}

func (m *DepthMatCache) Name() string { return "DepthMaterials" }

func (m *DepthMatCache) Instantiate(def *material.Def, callback func([]Material)) {
	if def == nil {
		def = &material.Def{}
	}

	// read vertex pointers from vertex format
	pointers := vertex.ParsePointers(def.VertexFormat)

	// fetch shader from cache
	shader := m.app.Shaders().Fetch(shader.Ref("depth"))

	// create material
	pipe := pipeline.New(
		m.app.Device(),
		pipeline.Args{
			Shader:     shader,
			Layout:     m.pipeLayout,
			Pass:       m.pass,
			Subpass:    MainSubpass,
			Pointers:   pointers,
			CullMode:   vertex.CullBack,
			DepthTest:  true,
			DepthWrite: true,
			DepthFunc:  core1_0.CompareOpLess,
			DepthClamp: def.DepthClamp,
			Primitive:  def.Primitive,
		})

	instances := make([]Material, m.frames)
	for i := range instances {
		desc := m.descLayout.Instantiate(m.app.Pool())
		instances[i] = &BasicMaterial{
			id:          def.Hash(),
			Pipeline:    pipe,
			Descriptors: desc,
			Objects:     NewObjectBuffer(desc.Objects.Size),
			Meshes:      m.app.Meshes(),
			Commands: command.NewIndirectIndexedDrawBuffer(m.app.Device(),
				fmt.Sprintf("DepthCommands:%d", i),
				desc.Objects.Size),
		}
	}

	callback(instances)
}

func (m *DepthMatCache) Destroy() {
	m.pipeLayout.Destroy()
	m.descLayout.Destroy()
}

func (m *DepthMatCache) Delete(mat []Material) {
	for _, m := range mat {
		m.Destroy()
	}
}
