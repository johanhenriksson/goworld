package descriptor

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Set interface {
	Ptr() core1_0.DescriptorSet
	Write(write core1_0.WriteDescriptorSet)
}

type set struct {
	device device.T
	layout SetLayout
	ptr    core1_0.DescriptorSet
}

func (s *set) Ptr() core1_0.DescriptorSet {
	return s.ptr
}

func (s *set) Write(write core1_0.WriteDescriptorSet) {
	write.DstSet = s.ptr
	if err := s.device.Ptr().UpdateDescriptorSets([]core1_0.WriteDescriptorSet{write}, nil); err != nil {
		panic(err)
	}
}
