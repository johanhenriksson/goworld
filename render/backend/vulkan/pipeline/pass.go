package pipeline

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	vk "github.com/vulkan-go/vulkan"
)

type Pass interface {
	device.Resource[vk.RenderPass]
}
