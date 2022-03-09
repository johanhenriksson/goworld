package descriptor

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"

	vk "github.com/vulkan-go/vulkan"
)

type Uniform struct {
	Binding int
	Stages  vk.ShaderStageFlags

	buffer buffer.T
	set    Set
}

var _ Descriptor = &Uniform{}

func (d *Uniform) Bind(set Set) {
	d.set = set
}

func (d *Uniform) Set(buffer buffer.T) {
	d.buffer = buffer
	d.Write()
}

func (d *Uniform) Write() {
	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstBinding:      uint32(d.Binding),
		DstArrayElement: 0,
		DescriptorCount: 1,
		DescriptorType:  vk.DescriptorTypeUniformBuffer,
		PBufferInfo: []vk.DescriptorBufferInfo{
			{
				Buffer: d.buffer.Ptr(),
				Offset: 0,
				Range:  vk.DeviceSize(vk.WholeSize),
			},
		},
	})
}

func (d *Uniform) LayoutBinding() vk.DescriptorSetLayoutBinding {
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(d.Binding),
		DescriptorType:  vk.DescriptorTypeUniformBuffer,
		DescriptorCount: 1,
		StageFlags:      d.Stages,
	}
}
