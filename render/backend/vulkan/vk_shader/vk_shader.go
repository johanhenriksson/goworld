package vk_shader

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex"

	vk "github.com/vulkan-go/vulkan"
)

type T[V any, K any, S any] interface {
	SetUniforms(frame int, uniforms []K)
	SetStorage(frame int, storage []S)
	Destroy()
	Bind(frame int, cmd command.Buffer)
	Attribute(name string) (shader.AttributeDesc, error)
	Layout() pipeline.Layout
}

type vk_shader[V any, U any, S any] struct {
	name       string
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
	attrs      shader.AttributeMap
}

type Args struct {
	Path       string
	Frames     int
	Bindings   []descriptor.Binding
	Pass       pipeline.Pass
	Attributes shader.AttributeMap
}

func New[V any, U any, S any](backend vulkan.T, args Args) T[V, U, S] {
	ubosize := 16 * 1024
	ssbosize := 1024 * 1024

	ubo := make([]buffer.T, args.Frames)
	ssbo := make([]buffer.T, args.Frames)
	for i := 0; i < args.Frames; i++ {
		ubo[i] = buffer.NewUniform(backend.Device(), ubosize)
		ssbo[i] = buffer.NewStorage(backend.Device(), ssbosize)
	}

	shaders := []pipeline.Shader{
		pipeline.NewShader(backend.Device(), fmt.Sprintf("assets/shaders/%s.vert.spv", args.Path), vk.ShaderStageVertexBit),
		pipeline.NewShader(backend.Device(), fmt.Sprintf("assets/shaders/%s.frag.spv", args.Path), vk.ShaderStageFragmentBit),
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

	dlayouts := make([]descriptor.T, 2*args.Frames)
	for i := 0; i < 2*args.Frames; i += args.Frames {
		dlayouts[i+0] = uboLayout
		dlayouts[i+1] = ssboLayout
	}

	// todo: calculate this based on input []bindings
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

	shader := &vk_shader[V, U, S]{
		name:       args.Path,
		frames:     args.Frames,
		backend:    backend,
		ubosize:    ubosize,
		shaders:    shaders,
		uboLayout:  uboLayout,
		ssboLayout: ssboLayout,
		dpool:      dpool,
		dsets:      dsets,
		layout:     layout,
		pass:       args.Pass,
		ubo:        ubo,
		ssbo:       ssbo,
		attrs:      args.Attributes,
	}

	for i := 0; i < args.Frames; i++ {
		shader.updateSets(i)
	}

	var vtx V
	pointers := vertex.ParsePointers(vtx)
	pointers.Bind(shader)

	fmt.Println(pointers)
	shader.pipe = pipeline.New(backend.Device(), nil, layout, args.Pass, shaders, pointers)

	return shader
}

func (s *vk_shader[V, U, S]) Name() string {
	return s.name
}

func (s *vk_shader[V, U, S]) Attribute(name string) (shader.AttributeDesc, error) {
	if attr, ok := s.attrs[name]; ok {
		return attr, nil
	}
	return shader.AttributeDesc{
		Name: name,
		Bind: -1,
	}, shader.ErrUnknownAttribute
}

func (s *vk_shader[V, U, S]) SetUniforms(frame int, uniforms []U) {
	s.ubo[frame].Write(uniforms, 0)
}

func (s *vk_shader[V, U, S]) SetStorage(frame int, storage []S) {
	s.ssbo[frame].Write(storage, 0)
}

func (s *vk_shader[V, U, S]) updateSets(frame int) {
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

func (s *vk_shader[V, U, S]) Bind(frame int, cmd command.Buffer) {
	cmd.CmdBindGraphicsPipeline(s.pipe)
	cmd.CmdBindGraphicsDescriptors(s.layout, s.dsets[2*frame:2*frame+2])
}

func (s *vk_shader[V, U, S]) Layout() pipeline.Layout {
	return s.layout
}

func (s *vk_shader[V, U, S]) Destroy() {
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
