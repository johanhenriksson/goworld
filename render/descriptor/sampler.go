package descriptor

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/texture"

	vk "github.com/vulkan-go/vulkan"
)

type Sampler struct {
	Stages vk.ShaderStageFlagBits

	binding int
	sampler vk.Sampler
	view    vk.ImageView
	set     Set
}

var _ Descriptor = &Sampler{}

func (d *Sampler) Initialize(device device.T) {}

func (d *Sampler) String() string {
	return fmt.Sprintf("Sampler:%d", d.binding)
}

func (d *Sampler) Destroy() {}

func (d *Sampler) Bind(set Set, binding int) {
	d.set = set
	d.binding = binding
}

func (d *Sampler) Set(tex texture.T) {
	d.sampler = tex.Ptr()
	d.view = tex.View().Ptr()
	d.write()
}

func (d *Sampler) LayoutBinding(binding int) vk.DescriptorSetLayoutBinding {
	d.binding = binding
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(binding),
		DescriptorType:  vk.DescriptorTypeCombinedImageSampler,
		DescriptorCount: 1,
		StageFlags:      vk.ShaderStageFlags(d.Stages),
	}
}

func (d *Sampler) BindingFlags() vk.DescriptorBindingFlags { return 0 }

func (d *Sampler) write() {
	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstSet:          d.set.Ptr(),
		DstBinding:      uint32(d.binding),
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