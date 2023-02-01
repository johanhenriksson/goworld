package descriptor

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/core1_2"
)

type SamplerArray struct {
	Count  int
	Stages core1_0.ShaderStageFlags

	binding int
	sampler []core1_0.Sampler
	view    []core1_0.ImageView
	set     Set
}

var _ Descriptor = &SamplerArray{}

func (d *SamplerArray) Initialize(device device.T) {
	if d.Count == 0 {
		panic("sampler array has count 0")
	}

	d.sampler = make([]core1_0.Sampler, d.Count)
	d.view = make([]core1_0.ImageView, d.Count)
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
		StageFlags:      core1_0.ShaderStageFlags(d.Stages),
	}
}

func (d *SamplerArray) BindingFlags() core1_2.DescriptorBindingFlags {
	return core1_2.DescriptorBindingFlags(
		core1_2.DescriptorBindingVariableDescriptorCount |
			core1_2.DescriptorBindingPartiallyBound |
			core1_2.DescriptorBindingUpdateAfterBind |
			core1_2.DescriptorBindingUpdateUnusedWhilePending)
}

func (d *SamplerArray) MaxCount() int {
	return d.Count
}

func (d *SamplerArray) Set(index int, tex texture.T) {
	d.sampler[index] = tex.Ptr()
	d.view[index] = tex.View().Ptr()
	d.write(index, 1)
}

func (d *SamplerArray) SetRange(textures []texture.T, offset int) {
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
