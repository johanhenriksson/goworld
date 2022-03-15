package descriptor

import vk "github.com/vulkan-go/vulkan"

type SetMock struct {
}

var _ Set = &SetMock{}

func (s *SetMock) Ptr() vk.DescriptorSet {
	return vk.NullDescriptorSet
}

func (s *SetMock) Write(write vk.WriteDescriptorSet) {
}
