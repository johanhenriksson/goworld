package descriptor

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type SamplerArray struct {
	Count  int
	Stages core1_0.ShaderStageFlags

	binding int
	set     Set

	// re-used update arrays
	info []core1_0.DescriptorImageInfo
}

var _ Descriptor = (*SamplerArray)(nil)
var _ VariableDescriptor = (*SamplerArray)(nil)

func (d *SamplerArray) Initialize(device *device.Device, set Set, binding int) {
	if d.Count == 0 {
		panic("sampler array has count 0")
	}

	d.set = set
	d.binding = binding

	d.info = make([]core1_0.DescriptorImageInfo, d.Count)
}

func (d *SamplerArray) String() string {
	return fmt.Sprintf("SamplerArray[%d]:%d", d.Count, d.binding)
}

func (d *SamplerArray) Destroy() {}

func (d *SamplerArray) LayoutBinding(binding int) core1_0.DescriptorSetLayoutBinding {
	return core1_0.DescriptorSetLayoutBinding{
		Binding:         binding,
		DescriptorType:  core1_0.DescriptorTypeCombinedImageSampler,
		DescriptorCount: d.Count,
		StageFlags:      d.Stages,
	}
}

func (d *SamplerArray) BindingFlags() ext_descriptor_indexing.DescriptorBindingFlags {
	return ext_descriptor_indexing.DescriptorBindingVariableDescriptorCount |
		ext_descriptor_indexing.DescriptorBindingPartiallyBound |
		ext_descriptor_indexing.DescriptorBindingUpdateAfterBind |
		ext_descriptor_indexing.DescriptorBindingUpdateUnusedWhilePending
}

func (d *SamplerArray) MaxCount() int {
	return d.Count
}

func (d *SamplerArray) Set(index int, tex *texture.Texture) {
	if index > d.Count {
		panic("out of bounds")
	}
	if tex == nil {
		panic("texture is null")
	}
	d.info[index] = core1_0.DescriptorImageInfo{
		Sampler:     tex.Ptr(),
		ImageView:   tex.View().Ptr(),
		ImageLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
	}
	d.write(index, 1)
}

func (d *SamplerArray) SetRange(offset int, textures texture.Array) {
	end := offset + len(textures)
	if end > d.Count {
		panic("out of bounds")
	}
	for i, tex := range textures {
		if tex == nil {
			panic(fmt.Sprintf("texture[%d] is null", i))
		}
		d.info[i] = core1_0.DescriptorImageInfo{
			Sampler:     tex.Ptr(),
			ImageView:   tex.View().Ptr(),
			ImageLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
		}
	}
	d.write(offset, len(textures))
}

func (d *SamplerArray) write(index, count int) {
	d.set.Write(core1_0.WriteDescriptorSet{
		DstSet:          d.set.Ptr(),
		DstBinding:      d.binding,
		DstArrayElement: index,
		DescriptorType:  core1_0.DescriptorTypeCombinedImageSampler,
		ImageInfo:       d.info[index : index+count],
	})
}
