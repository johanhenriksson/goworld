package descriptor

import (
	"github.com/vkngwrapper/core/v2/core1_0"
)

type SetMock struct {
}

var _ Set = &SetMock{}

func (s *SetMock) Ptr() core1_0.DescriptorSet {
	return nil
}

func (s *SetMock) Write(write core1_0.WriteDescriptorSet) {
}

func (s *SetMock) Destroy() {}

func (s *SetMock) adopt(Descriptor) {}
