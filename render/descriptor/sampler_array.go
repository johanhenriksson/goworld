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
	sampler []core1_0.Sampler
	view    []core1_0.ImageView
	set     Set

	// re-used update arrays
	info   []core1_0.DescriptorImageInfo
	writes []core1_0.WriteDescriptorSet
}

var _ Descriptor = &SamplerArray{}

func (d *SamplerArray) Initialize(device device.T) {
	if d.Count == 0 {
		panic("sampler array has count 0")
	}

	d.sampler = make([]core1_0.Sampler, d.Count)
	d.view = make([]core1_0.ImageView, d.Count)
	d.info = make([]core1_0.DescriptorImageInfo, 0, d.Count)
	d.writes = make([]core1_0.WriteDescriptorSet, 0, 100)
}

func (d *SamplerArray) String() string {
	return fmt.Sprintf("SamplerArray[%d]:%d", d.Count, d.binding)
}

func (d *SamplerArray) Destroy() {}

func (d *SamplerArray) Bind(set Set, binding int) {
	d.set = set
	d.binding = binding
}

func (d *SamplerArray) LayoutBinding(binding int) core1_0.DescriptorSetLayoutBinding {
	d.binding = binding
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

func (d *SamplerArray) Set(index int, tex texture.T) {
	if index > d.Count {
		panic("out of bounds")
	}
	if tex == nil {
		panic("texture is null")
	}
	d.sampler[index] = tex.Ptr()
	d.view[index] = tex.View().Ptr()
	d.write(index, 1)
}

func (d *SamplerArray) Clear(index int) {
	if index > d.Count {
		panic("out of bounds")
	}
	d.sampler[index] = nil
	d.view[index] = nil
	d.write(index, 1)
}

func (d *SamplerArray) SetRange(textures []texture.T, offset int) {
	end := offset + len(textures)
	if end > d.Count {
		panic("out of bounds")
	}
	for i, tex := range textures {
		if tex == nil {
			panic(fmt.Sprintf("texture[%d] is null", i))
		}
		d.sampler[offset+i] = tex.Ptr()
		d.view[offset+i] = tex.View().Ptr()
	}
	d.write(offset, len(textures))
}

func (d *SamplerArray) write(index, count int) {
	images := make([]core1_0.DescriptorImageInfo, count)
	for i := range images {
		images[i] = core1_0.DescriptorImageInfo{
			Sampler:     d.sampler[index+i],
			ImageView:   d.view[index+i],
			ImageLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
		}
	}

	d.set.Write(core1_0.WriteDescriptorSet{
		DstSet:          d.set.Ptr(),
		DstBinding:      d.binding,
		DstArrayElement: index,
		DescriptorType:  core1_0.DescriptorTypeCombinedImageSampler,
		ImageInfo:       images,
	})
}
