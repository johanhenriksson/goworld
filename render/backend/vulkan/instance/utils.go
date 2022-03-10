package instance

import (
	"fmt"

	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

func EnumerateExtensions() ([]string, error) {
	propCount := uint32(0)
	r := vk.EnumerateInstanceExtensionProperties("", &propCount, nil)
	if r != vk.Success {
		return nil, fmt.Errorf("get instance extensions count failed")
	}

	extensions := make([]vk.ExtensionProperties, propCount)
	r = vk.EnumerateInstanceExtensionProperties("", &propCount, extensions)
	if r != vk.Success {
		return nil, fmt.Errorf("enumerate instance extensions failed")
	}

	names := util.Map(extensions, func(ext vk.ExtensionProperties) string {
		ext.Deref()
		return vk.ToString(ext.ExtensionName[:])
	})

	return names, nil
}

func EnumerateInstanceLayers() ([]string, error) {
	count := uint32(0)
	r := vk.EnumerateInstanceLayerProperties(&count, nil)
	if r != vk.Success {
		return nil, fmt.Errorf("get instance extensions count failed")
	}

	layers := make([]vk.LayerProperties, count)
	r = vk.EnumerateInstanceLayerProperties(&count, layers)
	if r != vk.Success {
		return nil, fmt.Errorf("enumerate instance extensions failed")
	}

	names := util.Map(layers, func(layer vk.LayerProperties) string {
		layer.Deref()
		return vk.ToString(layer.LayerName[:])
	})

	return names, nil
}
