package descriptor

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"

	vk "github.com/vulkan-go/vulkan"
)

type Storage struct {
	Binding int
	Stages  vk.ShaderStageFlags

	buffer buffer.T
	set    Set
}

var _ Descriptor = &Storage{}

func (d *Storage) Bind(set Set) {
	d.set = set
}

func (d *Storage) Set(buffer buffer.T) {
	d.buffer = buffer
	d.Write()
}

func (d *Storage) LayoutBinding() vk.DescriptorSetLayoutBinding {
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(d.Binding),
		DescriptorType:  vk.DescriptorTypeStorageBuffer,
		DescriptorCount: 1,
		StageFlags:      d.Stages,
	}
}

func (d *Storage) Write() {
	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstSet:          d.set.Ptr(),
		DstBinding:      uint32(d.Binding),
		DstArrayElement: 0,
		DescriptorCount: 1,
		DescriptorType:  vk.DescriptorTypeStorageBuffer,
		PBufferInfo: []vk.DescriptorBufferInfo{
			{
				Buffer: d.buffer.Ptr(),
				Offset: 0,
				Range:  vk.DeviceSize(vk.WholeSize),
			},
		},
	})
}
