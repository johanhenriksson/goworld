package device

import (
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
	"github.com/vkngwrapper/extensions/v2/khr_maintenance3"
	"github.com/vkngwrapper/extensions/v2/khr_portability_subset"
	"github.com/vkngwrapper/extensions/v2/khr_storage_buffer_storage_class"
	"github.com/vkngwrapper/extensions/v2/khr_swapchain"
)

var deviceExtensions = []string{
	khr_swapchain.ExtensionName,
	khr_portability_subset.ExtensionName,
	khr_storage_buffer_storage_class.ExtensionName,
	ext_descriptor_indexing.ExtensionName,

	khr_maintenance3.ExtensionName,
}
