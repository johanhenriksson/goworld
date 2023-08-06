package material

import (
	"log"

	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/util"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// Materials combine pipelines and descriptors into a common unit.
type T[D descriptor.Set] interface {
	Destroy()
	TextureSlots() []texture.Slot
	Bind(cmd command.Buffer)
	Instantiate(descriptor.Pool) Instance[D]
	InstantiateMany(descriptor.Pool, int) []Instance[D]
}

type material[D descriptor.Set] struct {
	device  device.T
	dlayout descriptor.SetLayoutTyped[D]
	shader  shader.T
	layout  pipeline.Layout
	pipe    pipeline.T
	pass    renderpass.T
}

type Args struct {
	Shader     shader.T
	Pass       renderpass.T
	Subpass    renderpass.Name
	Pointers   vertex.Pointers
	Constants  []pipeline.PushConstant
	Primitive  vertex.Primitive
	DepthTest  bool
	DepthWrite bool
	DepthClamp bool
	DepthBias  float32
	DepthSlope float32
	DepthFunc  core1_0.CompareOp
	CullMode   vertex.CullMode
}

func New[D descriptor.Set](device device.T, args Args, descriptors D) T[D] {
	if device == nil {
		panic("device is nil")
	}
	if args.Shader == nil {
		panic("shader is nil")
	}

	for i, ptr := range args.Pointers {
		if index, kind, exists := args.Shader.Input(ptr.Name); exists {
			ptr.Bind(index, kind)
			args.Pointers[i] = ptr
		} else {
			log.Printf("no attribute in shader %s\n", ptr.Name)
		}
	}

	if args.Primitive == 0 {
		args.Primitive = vertex.Triangles
	}

	// create new descriptor set layout
	// ... this could be cached ...
	descLayout := descriptor.New(device, descriptors, args.Shader)

	// crete pipeline layout
	// ... this could be cached ...
	layout := pipeline.NewLayout(device, []descriptor.SetLayout{descLayout}, args.Constants)

	pipe := pipeline.New(device, pipeline.Args{
		Layout:   layout,
		Pass:     args.Pass,
		Subpass:  args.Subpass,
		Shader:   args.Shader,
		Pointers: args.Pointers,

		Primitive:  args.Primitive,
		DepthTest:  args.DepthTest,
		DepthWrite: args.DepthWrite,
		DepthClamp: args.DepthClamp,
		DepthFunc:  args.DepthFunc,
		CullMode:   args.CullMode,
	})

	return &material[D]{
		device: device,
		shader: args.Shader,

		dlayout: descLayout,
		layout:  layout,
		pipe:    pipe,
		pass:    args.Pass,
	}
}

func (m *material[D]) Bind(cmd command.Buffer) {
	cmd.CmdBindGraphicsPipeline(m.pipe)
}

func (m *material[D]) TextureSlots() []texture.Slot {
	return m.shader.Textures()
}

func (m *material[D]) Destroy() {
	m.dlayout.Destroy()
	m.pipe.Destroy()
	m.layout.Destroy()
}

func (m *material[D]) Instantiate(pool descriptor.Pool) Instance[D] {
	set := m.dlayout.Instantiate(pool)
	return &instance[D]{
		material: m,
		set:      set,
	}
}

func (m *material[D]) InstantiateMany(pool descriptor.Pool, n int) []Instance[D] {
	return util.Map(util.Range(0, n, 1), func(i int) Instance[D] { return m.Instantiate(pool) })
}
