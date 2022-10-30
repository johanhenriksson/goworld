package material

import (
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/pipeline"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Material templates allow for easy instantiation of pipelines for different vertex inputs, pipeline settings & shader modules
type Template[D descriptor.Set] interface {
	Name() string
	Instantiate(shader.T, vertex.Pointers) T[D]
	Destroy()
}

type template[D descriptor.Set] struct {
	device  device.T
	name    string
	dlayout descriptor.SetLayoutTyped[D]
	layout  pipeline.Layout
	shader  shader.T
	pass    renderpass.T
	subpass renderpass.Name
}

type TemplateArgs struct {
	Name      string
	Pass      renderpass.T
	Subpass   renderpass.Name
	Shader    shader.T
	Constants []pipeline.PushConstant
}

func NewTemplate[D descriptor.Set](device device.T, descriptors D, args TemplateArgs) Template[D] {
	// create new descriptor set layout
	descLayout := descriptor.New(device, descriptors, args.Shader)

	// crete pipeline layout
	layout := pipeline.NewLayout(device, []descriptor.SetLayout{descLayout}, args.Constants)

	return &template[D]{
		name:   args.Name,
		device: device,

		dlayout: descLayout,
		layout:  layout,
		shader:  args.Shader,
		pass:    args.Pass,
		subpass: args.Subpass,
	}
}

func (t *template[D]) Name() string {
	return t.name
}

func (t *template[D]) Destroy() {
	t.dlayout.Destroy()
	t.layout.Destroy()
}

func (t *template[D]) Instantiate(shader shader.T, pointers vertex.Pointers) T[D] {
	if shader == nil {
		shader = t.shader
	}

	pipe := pipeline.New(t.device, pipeline.Args{
		Layout:   t.layout,
		Pass:     t.pass,
		Subpass:  t.subpass,
		Shader:   shader,
		Pointers: pointers,

		Primitive:  vertex.Triangles,
		DepthTest:  true,
		DepthWrite: true,
	})

	return &material[D]{
		device: t.device,
		shader: shader,

		dlayout: t.dlayout,
		layout:  t.layout,
		pipe:    pipe,
		pass:    t.pass,
	}
}
