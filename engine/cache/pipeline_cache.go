package cache

import (
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Pipeline struct {
	ID       material.ID
	Slots    []texture.Slot
	pipeline *pipeline.Graphics
}

func (m *Pipeline) Bind(cmd *command.Buffer) {
	cmd.CmdBindGraphicsPipeline(m.pipeline)
}

func (m *Pipeline) Destroy() {
	m.pipeline.Destroy()
}

type PipelineCache T[*material.Def, *Pipeline]

type pipelineCache struct {
	device  *device.Device
	pass    *renderpass.Renderpass
	shaders ShaderCache
	layout  *pipeline.Layout
}

func NewPipelineCache(dev *device.Device, shaders ShaderCache, pass *renderpass.Renderpass, layout *pipeline.Layout) PipelineCache {
	return New[*material.Def, *Pipeline](&pipelineCache{
		device:  dev,
		pass:    pass,
		shaders: shaders,
		layout:  layout,
	})
}

func (m *pipelineCache) Name() string { return "ForwardMaterials" }

func (m *pipelineCache) Instantiate(def *material.Def, callback func(*Pipeline)) {
	if def == nil {
		def = material.StandardForward()
	}

	// read vertex pointers from vertex format
	pointers := vertex.ParsePointers(def.VertexFormat)

	// fetch shader from cache
	shader := m.shaders.Fetch(shader.Ref(def.Shader))

	// create material
	pipe := pipeline.New(
		m.device,
		pipeline.Args{
			Shader:     shader,
			Layout:     m.layout,
			Pass:       m.pass,
			Subpass:    "main",
			Pointers:   pointers,
			DepthTest:  def.DepthTest,
			DepthWrite: def.DepthWrite,
			DepthClamp: def.DepthClamp,
			DepthFunc:  def.DepthFunc,
			Primitive:  def.Primitive,
			CullMode:   def.CullMode,
		})

	callback(&Pipeline{
		ID:       def.Hash(),
		Slots:    shader.Textures(),
		pipeline: pipe,
	})
}

func (m *pipelineCache) Destroy() {
}

func (m *pipelineCache) Delete(mat *Pipeline) {
	mat.Destroy()
}
