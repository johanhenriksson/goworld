package shader

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"

	vk "github.com/vulkan-go/vulkan"
)

type T[K any] interface {
	SetUniforms(frame int, uniforms []K)
	Destroy()
	Bind(frame int, cmd command.Buffer)
}

type vk_shader[K any] struct {
	frames  int
	ubosize int
	backend vulkan.T
	ubo     []buffer.T
	shaders []pipeline.Shader
	dsets   []descriptor.Set
	dlayout descriptor.T
	dpool   descriptor.Pool
	layout  pipeline.Layout
	pipe    pipeline.T
	pass    pipeline.Pass
}

type Args struct {
	Path     string
	Frames   int
	Bindings []descriptor.Binding
	Pass     pipeline.Pass
}

func New[K any](backend vulkan.T, args Args) T[K] {
	ubosize := 3 * 16 * 4

	ubo := make([]buffer.T, args.Frames)
	for i := 0; i < args.Frames; i++ {
		ubo[i] = buffer.NewUniform(backend.Device(), ubosize)
	}

	shaders := []pipeline.Shader{
		pipeline.NewShader(backend.Device(), "assets/shaders/vk/color_f.vert.spv", vk.ShaderStageVertexBit),
		pipeline.NewShader(backend.Device(), "assets/shaders/vk/color_f.frag.spv", vk.ShaderStageFragmentBit),
	}

	dlayout := descriptor.New(backend.Device(), []descriptor.Binding{
		{
			Binding: 0,
			Type:    vk.DescriptorTypeUniformBuffer,
			Count:   1,
			Stages:  vk.ShaderStageFlags(vk.ShaderStageVertexBit),
		},
	})
	dlayouts := []descriptor.T{dlayout, dlayout}

	dpool := descriptor.NewPool(backend.Device(), []vk.DescriptorPoolSize{
		{
			Type:            vk.DescriptorTypeUniformBuffer,
			DescriptorCount: uint32(args.Frames),
		},
	})

	dsets := dpool.AllocateSets(dlayouts)

	layout := pipeline.NewLayout(backend.Device(), dlayouts)

	pipe := pipeline.New(backend.Device(), nil, layout, args.Pass, shaders)

	shader := &vk_shader[K]{
		backend: backend,
		ubosize: ubosize,
		shaders: shaders,
		dlayout: dlayout,
		dpool:   dpool,
		dsets:   dsets,
		layout:  layout,
		pipe:    pipe,
		pass:    args.Pass,
		ubo:     ubo,
	}

	for i := 0; i < args.Frames; i++ {
		shader.updateSets(i)
	}

	return shader
}

func (s *vk_shader[K]) SetUniforms(frame int, uniforms []K) {
	s.ubo[frame].Write(uniforms, 0)
}

func (s *vk_shader[K]) updateSets(frame int) {
	vk.UpdateDescriptorSets(s.backend.Device().Ptr(), 1, []vk.WriteDescriptorSet{
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          s.dsets[frame].Ptr(),
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
	}, 0, nil)
}

func (s *vk_shader[K]) Bind(frame int, cmd command.Buffer) {
	cmd.CmdBindGraphicsPipeline(s.pipe)
	cmd.CmdBindGraphicsDescriptors(s.layout, s.dsets[frame:frame+1])
}

func (s *vk_shader[K]) Destroy() {
	s.pipe.Destroy()
	s.layout.Destroy()
	s.dpool.Destroy()
	s.dlayout.Destroy()

	for _, shader := range s.shaders {
		shader.Destroy()
	}
	for _, ubo := range s.ubo {
		ubo.Destroy()
	}
}
