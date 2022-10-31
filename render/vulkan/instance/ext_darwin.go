package instance

import (
	vk "github.com/vulkan-go/vulkan"
)

var extensions = []string{
	vk.KhrSurfaceExtensionName,
	vk.ExtDebugReportExtensionName,
	"VK_EXT_metal_surface",
}
