package descriptor

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_texture"
	vk "github.com/vulkan-go/vulkan"
)

type SamplerArray struct {
	Count  int
	Stages vk.ShaderStageFlagBits

	binding int
	sampler []vk.Sampler
	view    []vk.ImageView
	set     Set
}

var _ Descriptor = &SamplerArray{}

func (d *SamplerArray) Initialize(device device.T) {
	if d.Count == 0 {
		panic("sampler array has count 0")
	}

	d.sampler = make([]vk.Sampler, d.Count)
	d.view = make([]vk.ImageView, d.Count)
}

func (d *SamplerArray) Destroy() {}

func (d *SamplerArray) Bind(set Set, binding int) {
	d.set = set
	d.binding = binding
}

func (d *SamplerArray) LayoutBinding(binding int) vk.DescriptorSetLayoutBinding {
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(binding),
		DescriptorType:  vk.DescriptorTypeCombinedImageSampler,
		DescriptorCount: uint32(d.Count),
		StageFlags:      vk.ShaderStageFlags(d.Stages),
	}
}

func (d *SamplerArray) BindingFlags() vk.DescriptorBindingFlags {
	return vk.DescriptorBindingFlags(
		vk.DescriptorBindingPartiallyBoundBit |
			vk.DescriptorBindingVariableDescriptorCountBit |
			vk.DescriptorBindingUpdateUnusedWhilePendingBit)
}

func (d *SamplerArray) MaxCount() int {
	return d.Count
}

func (d *SamplerArray) Set(index int, tex vk_texture.T) {
	d.sampler[index] = tex.Ptr()
	d.view[index] = tex.View().Ptr()
	d.write(index, 1)
}

func (d *SamplerArray) SetRange(textures []vk_texture.T, offset int) {
	end := offset + len(textures)
	if end >= d.Count {
		panic("out of bounds")
	}
	for i, tex := range textures {
		d.sampler[offset+i] = tex.Ptr()
		d.view[offset+i] = tex.View().Ptr()
	}
	d.write(offset, len(textures))
}

func (d *SamplerArray) write(index, count int) {
	images := make([]vk.DescriptorImageInfo, count)
	for i := range images {
		images[i] = vk.DescriptorImageInfo{
			Sampler:     d.sampler[i],
			ImageView:   d.view[i],
			ImageLayout: vk.ImageLayoutShaderReadOnlyOptimal,
		}
	}

	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstSet:          d.set.Ptr(),
		DstBinding:      uint32(d.binding),
		DstArrayElement: uint32(index),
		DescriptorCount: uint32(count),
		DescriptorType:  vk.DescriptorTypeCombinedImageSampler,
		PImageInfo:      images,
	})
}
