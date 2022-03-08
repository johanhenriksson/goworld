package vk_shader

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_texture"

	vk "github.com/vulkan-go/vulkan"
)

type DescriptorObject interface {
	Type() vk.DescriptorType
	Binding() descriptor.Binding
}

type Sampler interface {
	DescriptorObject
	SetTexture(backend vulkan.T, texture vk_texture.T, set descriptor.Set)
}

type StorageBuffer[T any] interface {
	Size() int
	Write(data T, offset int)
}

// {
// 	SType:           vk.StructureTypeWriteDescriptorSet,
// 	DstSet:          s.dsets[3*frame+1].Ptr(),
// 	DstBinding:      0,
// 	DstArrayElement: 0,
// 	DescriptorCount: 1,
// 	DescriptorType:  vk.DescriptorTypeStorageBuffer,
// 	PBufferInfo: []vk.DescriptorBufferInfo{
// 		{
// 			Buffer: s.ssbo[frame].Ptr(),
// 			Offset: 0,
// 			Range:  vk.DeviceSize(vk.WholeSize),
// 		},
// 	},
// },

type UniformBuffer[T any] interface {
	Size() int
	Write(data T, offset int)
}

// {
// 	SType:           vk.StructureTypeWriteDescriptorSet,
// 	DstSet:          s.dsets[3*frame+0].Ptr(),
// 	DstBinding:      0,
// 	DstArrayElement: 0,
// 	DescriptorCount: 1,
// 	DescriptorType:  vk.DescriptorTypeUniformBuffer,
// 	PBufferInfo: []vk.DescriptorBufferInfo{
// 		{
// 			Buffer: s.ubo[frame].Ptr(),
// 			Offset: 0,
// 			Range:  vk.DeviceSize(vk.WholeSize),
// 		},
// 	},
// },

type vk_sampler struct {
	bind    int
	stages  vk.ShaderStageFlags
	texture vk_texture.T
}

func NewSampler(binding int) Sampler {
	return &vk_sampler{
		bind:   binding,
		stages: vk.ShaderStageFlags(vk.ShaderStageAll),
	}
}

func (s *vk_sampler) Type() vk.DescriptorType {
	return vk.DescriptorTypeCombinedImageSampler
}

func (s *vk_sampler) Binding() descriptor.Binding {
	return descriptor.Binding{
		Type:    s.Type(),
		Stages:  s.stages,
		Binding: s.bind,
		Count:   1,
	}
}

func (s *vk_sampler) SetTexture(backend vulkan.T, texture vk_texture.T, set descriptor.Set) {
	s.texture = texture
	vk.UpdateDescriptorSets(backend.Device().Ptr(), 1, []vk.WriteDescriptorSet{
		{
			SType:           vk.StructureTypeWriteDescriptorSet,
			DstSet:          set.Ptr(),
			DstBinding:      uint32(s.bind),
			DstArrayElement: 0,
			DescriptorCount: 1,
			DescriptorType:  s.Type(),
			PImageInfo: []vk.DescriptorImageInfo{
				{
					Sampler:     texture.Ptr(),
					ImageView:   texture.View().Ptr(),
					ImageLayout: vk.ImageLayoutShaderReadOnlyOptimal,
				},
			},
		},
	}, 0, nil)
}
