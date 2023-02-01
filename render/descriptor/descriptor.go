package descriptor

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/core1_2"
)

type Descriptor interface {
	Initialize(device.T)
	LayoutBinding(int) core1_0.DescriptorSetLayoutBinding
	BindingFlags() core1_2.DescriptorBindingFlags
	Bind(Set, int)
	Destroy()
}

type VariableDescriptor interface {
	Descriptor
	MaxCount() int
}
