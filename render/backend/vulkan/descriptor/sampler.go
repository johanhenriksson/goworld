package descriptor

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/vk_texture"

	vk "github.com/vulkan-go/vulkan"
)

type Sampler struct {
	Binding int
	Stages  vk.ShaderStageFlags

	sampler vk.Sampler
	view    vk.ImageView
	layout  vk.ImageLayout
	set     Set
}

var _ Descriptor = &Sampler{}

func (d *Sampler) Bind(set Set) {
	d.set = set
}

func (d *Sampler) Set(tex vk_texture.T) {
	d.sampler = tex.Ptr()
	d.view = tex.View().Ptr()
	d.Write()
}

func (d *Sampler) LayoutBinding() vk.DescriptorSetLayoutBinding {
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(d.Binding),
		DescriptorType:  vk.DescriptorTypeCombinedImageSampler,
		DescriptorCount: 1,
		StageFlags:      d.Stages,
	}
}

func (d *Sampler) Write() {
	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstSet:          d.set.Ptr(),
		DstBinding:      uint32(d.Binding),
		DstArrayElement: 0,
		DescriptorCount: 1,
		DescriptorType:  vk.DescriptorTypeCombinedImageSampler,
		PImageInfo: []vk.DescriptorImageInfo{
			{
				Sampler:     d.sampler,
				ImageView:   d.view,
				ImageLayout: vk.ImageLayoutShaderReadOnlyOptimal,
			},
		},
	})
}