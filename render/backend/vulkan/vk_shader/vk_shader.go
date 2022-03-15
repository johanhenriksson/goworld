package vk_shader

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type Descriptors map[string]int

type T[V any, D descriptor.Set] interface {
	Destroy()
	Bind(frame int, cmd command.Buffer)
	Descriptors(frame int) D
	Layout() pipeline.Layout
}

type vk_shader[V any, D descriptor.Set] struct {
	name    string
	frames  int
	backend vulkan.T

	dsets   []D
	dpool   descriptor.Pool
	dlayout descriptor.SetLayoutTyped[D]

	shaders []pipeline.Shader
	layout  pipeline.Layout
	pipe    pipeline.T

	pass  renderpass.T
	attrs shader.AttributeMap
}

type Args struct {
	Path        string
	Frames      int
	Pass        renderpass.T
	Subpass     string
	Attributes  shader.AttributeMap
	Descriptors Descriptors
	Constants   []pipeline.PushConstant
}

func New[V any, D descriptor.Set](backend vulkan.T, descriptors D, args Args) T[V, D] {
	if args.Frames == 0 {
		args.Frames = backend.Frames()
	}

	shaders := []pipeline.Shader{
		pipeline.NewShader(backend.Device(), fmt.Sprintf("assets/shaders/%s.vert", args.Path), vk.ShaderStageVertexBit),
		pipeline.NewShader(backend.Device(), fmt.Sprintf("assets/shaders/%s.frag", args.Path), vk.ShaderStageFragmentBit),
	}

	// todo: put the descriptor pool somewhere else
	dpool := descriptor.NewPool(backend.Device(), []vk.DescriptorPoolSize{
		{
			Type:            vk.DescriptorTypeUniformBuffer,
			DescriptorCount: 100,
		},
		{
			Type:            vk.DescriptorTypeStorageBuffer,
			DescriptorCount: 100,
		},
		{
			Type:            vk.DescriptorTypeCombinedImageSampler,
			DescriptorCount: 100,
		},
		{
			Type:            vk.DescriptorTypeInputAttachment,
			DescriptorCount: 10,
		},
	})

	descLayout := descriptor.New(backend.Device(), descriptors)

	descSets := make([]D, args.Frames)
	for i := range descSets {
		dset := descLayout.Instantiate(dpool)
		descSets[i] = dset
	}

	layout := pipeline.NewLayout(backend.Device(), []descriptor.SetLayout{descLayout}, args.Constants)

	// todo: the pointers & pipeline stuff should be extracted into a material thing
	var vtx V
	pointers := vertex.ParsePointers(vtx)
	pointers.Bind(args.Attributes)

	pipe := pipeline.New(backend.Device(), pipeline.Args{
		Layout:   layout,
		Pass:     args.Pass,
		Subpass:  args.Subpass,
		Shaders:  shaders,
		Pointers: pointers,

		Primitive:  vertex.Triangles,
		DepthTest:  true,
		DepthWrite: true,
	})

	return &vk_shader[V, D]{
		name:    args.Path,
		backend: backend,
		shaders: shaders,
		frames:  args.Frames,

		dsets:   descSets,
		dlayout: descLayout,
		dpool:   dpool,

		attrs:  args.Attributes,
		layout: layout,
		pipe:   pipe,
		pass:   args.Pass,
	}
}

func (s *vk_shader[V, D]) Name() string {
	return s.name
}

func (s *vk_shader[V, D]) Descriptors(frame int) D {
	return s.dsets[frame%s.frames]
}

func (s *vk_shader[V, D]) Layout() pipeline.Layout {
	return s.layout
}

func (s *vk_shader[V, D]) Bind(frame int, cmd command.Buffer) {
	cmd.CmdBindGraphicsPipeline(s.pipe)
	cmd.CmdBindGraphicsDescriptor(s.layout, s.dsets[frame%s.frames])
}

func (s *vk_shader[V, D]) Destroy() {
	s.dlayout.Destroy()
	s.dpool.Destroy()

	s.pipe.Destroy()
	s.layout.Destroy()

	for _, shader := range s.shaders {
		shader.Destroy()
	}
}
