package descriptor

import (
	vk "github.com/vulkan-go/vulkan"
)

type Set interface {
	Ptr() vk.DescriptorSet
}

type set struct {
	ptr vk.DescriptorSet
}

func (s *set) Ptr() vk.DescriptorSet {
	return s.ptr
}
