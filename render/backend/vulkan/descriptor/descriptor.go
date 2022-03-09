package descriptor

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Descriptor interface {
	Initialize(device.T)
	LayoutBinding() vk.DescriptorSetLayoutBinding
	Write()
	Bind(Set)
	Destroy()
}
