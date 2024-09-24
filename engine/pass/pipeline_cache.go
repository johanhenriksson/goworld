package pass

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Pipeline struct {
	id       material.ID
	slots    []texture.Slot
	pipeline *pipeline.Pipeline
}

func (m *Pipeline) Bind(cmd *command.Buffer) {
	cmd.CmdBindGraphicsPipeline(m.pipeline)
}

func (m *Pipeline) Destroy() {
	m.pipeline.Destroy()
}

type PipelineCache struct {
	app     engine.App
	pass    *renderpass.Renderpass
	shaders cache.ShaderCache
	frames  int
	layout  *pipeline.Layout
}

func NewPipelineCache(app engine.App, pass *renderpass.Renderpass, frames int, layout *pipeline.Layout) cache.T[*material.Def, *Pipeline] {
	return cache.New[*material.Def, *Pipeline](&PipelineCache{
		app:     app,
		pass:    pass,
		shaders: app.Shaders(),
		frames:  frames,
		layout:  layout,
	})
}

func (m *PipelineCache) Name() string { return "ForwardMaterials" }

func (m *PipelineCache) Instantiate(def *material.Def, callback func(*Pipeline)) {
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
			Layout:     m.layout,
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

	callback(&Pipeline{
		id:       def.Hash(),
		slots:    shader.Textures(),
		pipeline: pipe,
	})
}

func (m *PipelineCache) Destroy() {
}

func (m *PipelineCache) Delete(mat *Pipeline) {
	mat.Destroy()
}
