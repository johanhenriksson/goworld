package descriptor

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Set interface {
	Ptr() core1_0.DescriptorSet
	Write(write core1_0.WriteDescriptorSet)
	Destroy()

	adopt(Descriptor)
}

type set struct {
	device *device.Device
	layout SetLayout
	ptr    core1_0.DescriptorSet
	items  []Descriptor
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

func (s *set) Destroy() {
	for _, d := range s.items {
		d.Destroy()
	}
}

// adopt a descriptor, freeing it when the set is destroyed
func (s *set) adopt(d Descriptor) {
	s.items = append(s.items, d)
}
