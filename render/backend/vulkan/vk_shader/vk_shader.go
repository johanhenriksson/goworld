package vk_shader

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/renderpass"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_texture"
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

	SetTexture(frame int, name string, tex vk_texture.T)
}

type SamplerMap map[string]int

type vk_shader[V any, U any, S any] struct {
	name       string
	frames     int
	ubosize    int
	backend    vulkan.T
	ubo        []buffer.T
	ssbo       []buffer.T
	shaders    []pipeline.Shader
	dsets      []descriptor.Set
	sets       int
	uboLayout  descriptor.T
	uboSet     int
	ssboLayout descriptor.T
	ssboSet    int
	texLayout  descriptor.T
	texSet     int
	dpool      descriptor.Pool
	layout     pipeline.Layout
	pipe       pipeline.T
	pass       renderpass.T
	attrs      shader.AttributeMap
	samplers   map[string]Sampler
}

type Args struct {
	Path       string
	Bindings   []descriptor.Binding
	Pass       renderpass.T
	Attributes shader.AttributeMap
	Samplers   SamplerMap
}

func New[V any, U any, S any](backend vulkan.T, args Args) T[V, U, S] {
	ubosize := 16 * 1024
	ssbosize := 1024 * 1024
	frames := backend.Frames()

	ubo := make([]buffer.T, frames)
	ssbo := make([]buffer.T, frames)
	for i := 0; i < frames; i++ {
		ubo[i] = buffer.NewUniform(backend.Device(), ubosize)
		ssbo[i] = buffer.NewStorage(backend.Device(), ssbosize)
	}

	shaders := []pipeline.Shader{
		pipeline.NewShader(backend.Device(), fmt.Sprintf("assets/shaders/%s.vert.spv", args.Path), vk.ShaderStageVertexBit),
		pipeline.NewShader(backend.Device(), fmt.Sprintf("assets/shaders/%s.frag.spv", args.Path), vk.ShaderStageFragmentBit),
	}

	var texLayout descriptor.T

	sets := 2
	uboSet, ssboSet, texSet := -1, -1, -1
	poolsizes := make([]vk.DescriptorPoolSize, 0, 3)

	uboSet = 0
	uboLayout := descriptor.New(backend.Device(), []descriptor.Binding{
		{
			Binding: 0,
			Type:    vk.DescriptorTypeUniformBuffer,
			Count:   1,
			Stages:  vk.ShaderStageFlags(vk.ShaderStageAll),
		},
	})
	poolsizes = append(poolsizes, vk.DescriptorPoolSize{
		Type:            vk.DescriptorTypeUniformBuffer,
		DescriptorCount: uint32(frames),
	})

	ssboSet = 1
	ssboLayout := descriptor.New(backend.Device(), []descriptor.Binding{
		{
			Binding: 0,
			Type:    vk.DescriptorTypeStorageBuffer,
			Count:   1,
			Stages:  vk.ShaderStageFlags(vk.ShaderStageAll),
		},
	})
	poolsizes = append(poolsizes, vk.DescriptorPoolSize{
		Type:            vk.DescriptorTypeStorageBuffer,
		DescriptorCount: uint32(frames),
	})

	texBinds := make([]descriptor.Binding, 0, len(args.Samplers))
	samplers := make(map[string]Sampler)
	for name, binding := range args.Samplers {
		sampler := NewSampler(binding)
		texBinds = append(texBinds, sampler.Binding())
		samplers[name] = sampler
	}
	if len(texBinds) > 0 {
		texLayout = descriptor.New(backend.Device(), texBinds)
		poolsizes = append(poolsizes, vk.DescriptorPoolSize{
			Type:            vk.DescriptorTypeCombinedImageSampler,
			DescriptorCount: uint32(frames * len(texBinds)),
		})
		texSet = sets
		sets++
	}

	dlayouts := make([]descriptor.T, sets*frames)
	for i := 0; i < frames; i++ {
		if uboSet >= 0 {
			dlayouts[sets*i+uboSet] = uboLayout
		}
		if ssboSet >= 0 {
			dlayouts[sets*i+ssboSet] = ssboLayout
		}
		if texSet >= 0 {
			dlayouts[sets*i+texSet] = texLayout
		}
	}

	log.Println("descriptor sets", sets)

	// todo: calculate this based on input []bindings
	dpool := descriptor.NewPool(backend.Device(), poolsizes)
	dsets := dpool.AllocateSets(dlayouts)

	layout := pipeline.NewLayout(backend.Device(), dlayouts)

	shader := &vk_shader[V, U, S]{
		name:    args.Path,
		frames:  frames,
		backend: backend,
		ubosize: ubosize,
		shaders: shaders,

		sets:       sets,
		uboSet:     uboSet,
		uboLayout:  uboLayout,
		ssboSet:    ssboSet,
		ssboLayout: ssboLayout,
		texSet:     texSet,
		texLayout:  texLayout,

		dpool:    dpool,
		dsets:    dsets,
		layout:   layout,
		pass:     args.Pass,
		ubo:      ubo,
		ssbo:     ssbo,
		attrs:    args.Attributes,
		samplers: samplers,
	}

	for i := 0; i < frames; i++ {
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

func (s *vk_shader[V, U, S]) SetTexture(frame int, name string, tex vk_texture.T) {
	sampler := s.samplers[name]
	sampler.SetTexture(s.backend, tex, s.dsets[s.sets*frame+s.texSet])
}

func (s *vk_shader[V, U, S]) updateSets(frame int) {
	vk.UpdateDescriptorSets(s.backend.Device().Ptr(), 2, []vk.WriteDescriptorSet{
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          s.dsets[s.sets*frame+0].Ptr(),
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
			DstSet:          s.dsets[s.sets*frame+1].Ptr(),
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
	cmd.CmdBindGraphicsDescriptors(s.layout, s.dsets[s.sets*frame:s.sets*(frame+1)])
}

func (s *vk_shader[V, U, S]) Layout() pipeline.Layout {
	return s.layout
}

func (s *vk_shader[V, U, S]) Destroy() {
	s.pipe.Destroy()
	s.layout.Destroy()
	s.dpool.Destroy()

	if s.uboLayout != nil {
		s.uboLayout.Destroy()
		s.uboLayout = nil
	}

	if s.ssboLayout != nil {
		s.ssboLayout.Destroy()
		s.ssboLayout = nil
	}

	if s.texLayout != nil {
		s.texLayout.Destroy()
		s.texLayout = nil
	}

	for _, shader := range s.shaders {
		shader.Destroy()
	}

	for i := 0; i < s.frames; i++ {
		s.ubo[i].Destroy()
		s.ssbo[i].Destroy()
	}
}
