package vk_shader

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Descriptors map[string]int

type T[D descriptor.Set] interface {
	Destroy()
	Bind(frame int, cmd command.Buffer)
	Descriptors(frame int) D
	Layout() pipeline.Layout
}

type vk_shader[D descriptor.Set] struct {
	name    string
	frames  int
	backend vulkan.T

	dsets   []D
	dlayout descriptor.SetLayoutTyped[D]

	shader pipeline.Shader
	layout pipeline.Layout
	pipe   pipeline.T

	pass  renderpass.T
	attrs shader.AttributeMap
}

type Args struct {
	Path       string
	Frames     int
	Pass       renderpass.T
	Subpass    string
	Attributes shader.AttributeMap
	Pointers   vertex.Pointers
	Constants  []pipeline.PushConstant
}

func New[D descriptor.Set](backend vulkan.T, args Args, descriptors D) T[D] {
	if args.Frames == 0 {
		args.Frames = backend.Frames()
	}

	// instantiate shader modules
	// ... this could be cached ...
	shader := pipeline.NewShader(backend.Device(), args.Path)

	// create new descriptor set layout
	// ... this could be cached ...
	descLayout := descriptor.New(backend.Device(), descriptors)

	// crete pipeline layout
	// ... this could be cached ...
	layout := pipeline.NewLayout(backend.Device(), []descriptor.SetLayout{descLayout}, args.Constants)

	// instantiate one descriptor set per frame
	descSets := make([]D, args.Frames)
	for i := range descSets {
		dset := descLayout.Instantiate(descriptor.GlobalPool)
		descSets[i] = dset
	}

	// todo: the pointers & pipeline stuff should be extracted into a material thing
	args.Pointers.Bind(args.Attributes)

	pipe := pipeline.New(backend.Device(), pipeline.Args{
		Layout:   layout,
		Pass:     args.Pass,
		Subpass:  args.Subpass,
		Shader:   shader,
		Pointers: args.Pointers,

		Primitive:  vertex.Triangles,
		DepthTest:  true,
		DepthWrite: true,
	})

	return &vk_shader[D]{
		name:    args.Path,
		backend: backend,
		shader:  shader,
		frames:  args.Frames,

		dsets:   descSets,
		dlayout: descLayout,

		attrs:  args.Attributes,
		layout: layout,
		pipe:   pipe,
		pass:   args.Pass,
	}
}

func (s *vk_shader[D]) Name() string {
	return s.name
}

func (s *vk_shader[D]) Descriptors(frame int) D {
	return s.dsets[frame%s.frames]
}

func (s *vk_shader[D]) Layout() pipeline.Layout {
	return s.layout
}

func (s *vk_shader[D]) Bind(frame int, cmd command.Buffer) {
	cmd.CmdBindGraphicsPipeline(s.pipe)
	cmd.CmdBindGraphicsDescriptor(s.layout, s.dsets[frame%s.frames])
}

func (s *vk_shader[D]) Destroy() {
	s.dlayout.Destroy()
	s.pipe.Destroy()
	s.layout.Destroy()
	s.shader.Destroy()
}
