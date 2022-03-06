package shader

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"

	vk "github.com/vulkan-go/vulkan"
)

type T[K any, S any] interface {
	SetUniforms(frame int, uniforms []K)
	SetStorage(frame int, storage []S)
	Destroy()
	Bind(frame int, cmd command.Buffer)
	Layout() pipeline.Layout
}

type vk_shader[K any, S any] struct {
	frames     int
	ubosize    int
	backend    vulkan.T
	ubo        []buffer.T
	ssbo       []buffer.T
	shaders    []pipeline.Shader
	dsets      []descriptor.Set
	uboLayout  descriptor.T
	ssboLayout descriptor.T
	dpool      descriptor.Pool
	layout     pipeline.Layout
	pipe       pipeline.T
	pass       pipeline.Pass
}

type Args struct {
	Path     string
	Frames   int
	Bindings []descriptor.Binding
	Pass     pipeline.Pass
}

func New[K any, S any](backend vulkan.T, args Args) T[K, S] {
	ubosize := 16 * 1024
	ssbosize := 1024 * 1024

	ubo := make([]buffer.T, args.Frames)
	ssbo := make([]buffer.T, args.Frames)
	for i := 0; i < args.Frames; i++ {
		ubo[i] = buffer.NewUniform(backend.Device(), ubosize)
		ssbo[i] = buffer.NewStorage(backend.Device(), ssbosize)
	}

	shaders := []pipeline.Shader{
		pipeline.NewShader(backend.Device(), "assets/shaders/vk/color_f.vert.spv", vk.ShaderStageVertexBit),
		pipeline.NewShader(backend.Device(), "assets/shaders/vk/color_f.frag.spv", vk.ShaderStageFragmentBit),
	}

	uboLayout := descriptor.New(backend.Device(), []descriptor.Binding{
		{
			Binding: 0,
			Type:    vk.DescriptorTypeUniformBuffer,
			Count:   1,
			Stages:  vk.ShaderStageFlags(vk.ShaderStageAll),
		},
	})
	ssboLayout := descriptor.New(backend.Device(), []descriptor.Binding{
		{
			Binding: 0,
			Type:    vk.DescriptorTypeStorageBuffer,
			Count:   1,
			Stages:  vk.ShaderStageFlags(vk.ShaderStageAll),
		},
	})
	dlayouts := []descriptor.T{uboLayout, ssboLayout, uboLayout, ssboLayout}

	dpool := descriptor.NewPool(backend.Device(), []vk.DescriptorPoolSize{
		{
			Type:            vk.DescriptorTypeUniformBuffer,
			DescriptorCount: uint32(args.Frames),
		},
		{
			Type:            vk.DescriptorTypeStorageBuffer,
			DescriptorCount: uint32(args.Frames),
		},
	})

	dsets := dpool.AllocateSets(dlayouts)

	layout := pipeline.NewLayout(backend.Device(), dlayouts)

	pipe := pipeline.New(backend.Device(), nil, layout, args.Pass, shaders)

	shader := &vk_shader[K, S]{
		frames:     args.Frames,
		backend:    backend,
		ubosize:    ubosize,
		shaders:    shaders,
		uboLayout:  uboLayout,
		ssboLayout: ssboLayout,
		dpool:      dpool,
		dsets:      dsets,
		layout:     layout,
		pipe:       pipe,
		pass:       args.Pass,
		ubo:        ubo,
		ssbo:       ssbo,
	}

	for i := 0; i < args.Frames; i++ {
		shader.updateSets(i)
	}

	return shader
}

func (s *vk_shader[K, S]) SetUniforms(frame int, uniforms []K) {
	s.ubo[frame].Write(uniforms, 0)
}

func (s *vk_shader[K, S]) SetStorage(frame int, storage []S) {
	s.ssbo[frame].Write(storage, 0)
}

func (s *vk_shader[K, S]) updateSets(frame int) {
	vk.UpdateDescriptorSets(s.backend.Device().Ptr(), 2, []vk.WriteDescriptorSet{
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          s.dsets[2*frame+0].Ptr(),
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorCount: 1,
			DescriptorType:  vk.DescriptorTypeUniformBuffer,
			PBufferInfo: []vk.DescriptorBufferInfo{
				{
					Buffer: s.ubo[frame].Ptr(),
					Offset: 0,
					Range:  vk.DeviceSize(vk.WholeSize),
				},
			},
		},
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          s.dsets[2*frame+1].Ptr(),
			DstBinding:      0,
			DstArrayElement: 0,
			DescriptorCount: 1,
			DescriptorType:  vk.DescriptorTypeStorageBuffer,
			PBufferInfo: []vk.DescriptorBufferInfo{
				{
					Buffer: s.ssbo[frame].Ptr(),
					Offset: 0,
					Range:  vk.DeviceSize(vk.WholeSize),
				},
			},
		},
	}, 0, nil)
}

func (s *vk_shader[K, S]) Bind(frame int, cmd command.Buffer) {
	cmd.CmdBindGraphicsPipeline(s.pipe)
	cmd.CmdBindGraphicsDescriptors(s.layout, s.dsets[2*frame:2*frame+2])
}

func (s *vk_shader[K, S]) Layout() pipeline.Layout {
	return s.layout
}

func (s *vk_shader[K, S]) Destroy() {
	s.pipe.Destroy()
	s.layout.Destroy()
	s.dpool.Destroy()
	s.uboLayout.Destroy()
	s.ssboLayout.Destroy()

	for _, shader := range s.shaders {
		shader.Destroy()
	}
	for i := 0; i < s.frames; i++ {
		s.ubo[i].Destroy()
		s.ssbo[i].Destroy()
	}
}
