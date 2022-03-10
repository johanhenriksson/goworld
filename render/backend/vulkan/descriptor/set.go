package descriptor

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	vk "github.com/vulkan-go/vulkan"
)

type Set interface {
	Ptr() vk.DescriptorSet
	Write(write vk.WriteDescriptorSet)
}

type set struct {
	device device.T
	layout SetLayout
	ptr    vk.DescriptorSet
}

func (s *set) Ptr() vk.DescriptorSet {
	return s.ptr
}

func (s *set) Layout() SetLayout {
	return s.layout
}

func (s *set) Write(write vk.WriteDescriptorSet) {
	write.SType = vk.StructureTypeWriteDescriptorSet
	write.DstSet = s.ptr
	vk.UpdateDescriptorSets(s.device.Ptr(), 1, []vk.WriteDescriptorSet{write}, 0, nil)
}
