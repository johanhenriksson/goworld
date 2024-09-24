package instance

import (
	"github.com/vkngwrapper/extensions/v2/ext_debug_utils"
	"github.com/vkngwrapper/extensions/v2/khr_get_physical_device_properties2"
	"github.com/vkngwrapper/extensions/v2/khr_portability_enumeration"
	"github.com/vkngwrapper/extensions/v2/khr_surface"
)

var extensions = []string{
	khr_surface.ExtensionName,
	ext_debug_utils.ExtensionName,
	khr_get_physical_device_properties2.ExtensionName,
	khr_portability_enumeration.ExtensionName,

	"VK_EXT_debug_report",
	"VK_EXT_metal_surface",
}
