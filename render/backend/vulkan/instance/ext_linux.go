package instance

import (
	vk "github.com/vulkan-go/vulkan"
)

var extensions = []string{
	vk.KhrSurfaceExtensionName,
	// vk.ExtDebugReportExtensionName,
	"VK_KHR_xcb_surface",
}
