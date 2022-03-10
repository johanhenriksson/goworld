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
}

type vk_shader[V any, D descriptor.Set] struct {
	name    string
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
	Pass        renderpass.T
	Attributes  shader.AttributeMap
	Descriptors Descriptors
}

func New[V any, D descriptor.Set](backend vulkan.T, descriptors D, args Args) T[V, D] {
	shaders := []pipeline.Shader{
		pipeline.NewShader(backend.Device(), fmt.Sprintf("assets/shaders/%s.vert.spv", args.Path), vk.ShaderStageVertexBit),
		pipeline.NewShader(backend.Device(), fmt.Sprintf("assets/shaders/%s.frag.spv", args.Path), vk.ShaderStageFragmentBit),
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
	})

	descLayout := descriptor.New(backend.Device(), dpool, descriptors)

	descSets := make([]D, backend.Frames())
	for i := range descSets {
		dset := descLayout.Allocate()
		descSets[i] = dset
	}

	layout := pipeline.NewLayout(backend.Device(), []descriptor.SetLayout{descLayout})

	// todo: the pointers & pipeline stuff should be extracted into a material thing
	var vtx V
	pointers := vertex.ParsePointers(vtx)
	pointers.Bind(args.Attributes)

	pipe := pipeline.New(backend.Device(), pipeline.Args{
		Layout:   layout,
		Pass:     args.Pass,
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
	return s.dsets[frame]
}

func (s *vk_shader[V, D]) Bind(frame int, cmd command.Buffer) {
	cmd.CmdBindGraphicsPipeline(s.pipe)
	cmd.CmdBindGraphicsDescriptor(s.layout, s.dsets[frame])
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
