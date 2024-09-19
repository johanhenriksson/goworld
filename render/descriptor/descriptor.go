package descriptor

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type Descriptor interface {
	// Initialize the descriptor with the provided device, set, and binding index.
	// This is called automatically during Set instantiation.
	Initialize(dev *device.Device, set Set, binding int)

	BindingFlags() ext_descriptor_indexing.DescriptorBindingFlags

	// Destroy the descriptor.
	// Releasing any resources it owns, such as uniform buffers.
	Destroy()

	// LayoutBinding returns a descriptor set layout binding for this descriptor
	// using the provided binding index. This is useful for creating descriptor
	// set layouts including this descriptor.
	LayoutBinding(int) core1_0.DescriptorSetLayoutBinding
}

type VariableDescriptor interface {
	Descriptor
	MaxCount() int
}
