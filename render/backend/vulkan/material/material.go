package material

import (
	"log"

	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
	"github.com/johanhenriksson/goworld/util"
)

type T[D descriptor.Set] interface {
	Destroy()
	Bind(cmd command.Buffer)
	Layout() pipeline.Layout
	Instantiate() Instance[D]
	InstantiateMany(int) []Instance[D]
}

type material[D descriptor.Set] struct {
	backend vulkan.T
	dlayout descriptor.SetLayoutTyped[D]
	shader  shader.T
	layout  pipeline.Layout
	pipe    pipeline.T
	pass    renderpass.T
}

type Args struct {
	Shader    shader.T
	Pass      renderpass.T
	Subpass   string
	Pointers  vertex.Pointers
	Constants []pipeline.PushConstant
}

func New[D descriptor.Set](backend vulkan.T, args Args, descriptors D) T[D] {
	// instantiate shader modules
	// ... this could be cached ...

	for i, ptr := range args.Pointers {
		if index, kind, exists := args.Shader.Input(ptr.Name); exists {
			ptr.Bind(index, kind)
			args.Pointers[i] = ptr
		} else {
			log.Printf("no attribute in shader %s\n", ptr.Name)
		}
	}

	// create new descriptor set layout
	// ... this could be cached ...
	descLayout := descriptor.New(backend.Device(), descriptors, args.Shader)

	// crete pipeline layout
	// ... this could be cached ...
	layout := pipeline.NewLayout(backend.Device(), []descriptor.SetLayout{descLayout}, args.Constants)

	pipe := pipeline.New(backend.Device(), pipeline.Args{
		Layout:   layout,
		Pass:     args.Pass,
		Subpass:  args.Subpass,
		Shader:   args.Shader,
		Pointers: args.Pointers,

		Primitive:  vertex.Triangles,
		DepthTest:  true,
		DepthWrite: true,
	})

	return &material[D]{
		backend: backend,
		shader:  args.Shader,

		dlayout: descLayout,
		layout:  layout,
		pipe:    pipe,
		pass:    args.Pass,
	}
}

func (m *material[D]) Layout() pipeline.Layout {
	return m.layout
}

func (m *material[D]) Bind(cmd command.Buffer) {
	cmd.CmdBindGraphicsPipeline(m.pipe)
}

func (m *material[D]) Destroy() {
	m.dlayout.Destroy()
	m.pipe.Destroy()
	m.layout.Destroy()
	m.shader.Destroy()
}

func (m *material[D]) Instantiate() Instance[D] {
	set := m.dlayout.Instantiate(descriptor.GlobalPool)
	return &instance[D]{
		material: m,
		set:      set,
	}
}

func (m *material[D]) InstantiateMany(n int) []Instance[D] {
	return util.Map(util.Range(0, n, 1), func(i int) Instance[D] { return m.Instantiate() })
}
