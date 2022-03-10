package descriptor

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Descriptor interface {
	Initialize(device.T)
	LayoutBinding() vk.DescriptorSetLayoutBinding
	BindingFlags() vk.DescriptorBindingFlags
	Bind(Set)
	Destroy()
}

type VariableDescriptor interface {
	Descriptor
	MaxCount() int
}
