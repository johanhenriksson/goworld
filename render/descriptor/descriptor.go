package descriptor

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type Descriptor interface {
	Initialize(device.T)
	LayoutBinding(int) core1_0.DescriptorSetLayoutBinding
	BindingFlags() ext_descriptor_indexing.DescriptorBindingFlags
	Bind(Set, int)
	Destroy()
}

type VariableDescriptor interface {
	Descriptor
	MaxCount() int
}
