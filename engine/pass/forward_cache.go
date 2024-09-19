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
	app    engine.App
	pass   *renderpass.Renderpass
	layout *descriptor.Layout[*ForwardDescriptors]
	lookup ShadowmapLookupFn
	frames int
}

func NewForwardMaterialCache(app engine.App, pass *renderpass.Renderpass, frames int, lookup ShadowmapLookupFn) MaterialCache {
	layout := descriptor.NewLayout(app.Device(), "Forward", &ForwardDescriptors{
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
	})
	return cache.New[*material.Def, []Material](&ForwardMatCache{
		app:    app,
		pass:   pass,
		lookup: lookup,
		frames: frames,
		layout: layout,
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
	shader := m.app.Shaders().Fetch(shader.Ref(def.Shader))

	// create material
	pipe := pipeline.New(
		m.app.Device(),
		pipeline.Args{
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
		m.layout)

	instances := make([]Material, m.frames)
	for i := range instances {
		desc := m.layout.Instantiate(m.app.Pool())
		textures := cache.NewSamplerCache(m.app.Textures(), desc.Textures)

		instances[i] = &ForwardMaterial{
			id:          def.Hash(),
			Pipeline:    pipe,
			Descriptors: desc,
			Objects:     NewObjectBuffer(desc.Objects.Size),
			Lights:      NewLightBuffer(desc.Lights.Size),
			Shadows:     NewShadowCache(textures, m.lookup),
			Textures:    textures,
			Meshes:      m.app.Meshes(),
			Commands: command.NewIndirectDrawBuffer(m.app.Device(),
				fmt.Sprintf("ForwardCommands:%d", i),
				desc.Objects.Size),
		}
	}

	callback(instances)
}

func (m *ForwardMatCache) Destroy() {
	m.layout.Destroy()
}

func (m *ForwardMatCache) Delete(mat []Material) {
	for _, m := range mat {
		m.Destroy()
	}
}
