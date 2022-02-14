package descriptor

import (
	vk "github.com/vulkan-go/vulkan"
)

type Set interface {
	Ptr() vk.DescriptorSet
	Layout() T
}

type set struct {
	layout T
	ptr    vk.DescriptorSet
}

func (s *set) Ptr() vk.DescriptorSet {
	return s.ptr
}

func (s *set) Layout() T {
	return s.layout
}
