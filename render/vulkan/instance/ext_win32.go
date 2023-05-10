package instance

import (
	"github.com/vkngwrapper/extensions/v2/ext_debug_utils"
	"github.com/vkngwrapper/extensions/v2/khr_get_physical_device_properties2"
	"github.com/vkngwrapper/extensions/v2/khr_surface"
)

var extensions = []string{
	khr_surface.ExtensionName,
	khr_get_physical_device_properties2.ExtensionName,
	ext_debug_utils.ExtensionName,

	"VK_KHR_win32_surface",
}
