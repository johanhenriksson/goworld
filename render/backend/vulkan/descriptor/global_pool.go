package descriptor

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	vk "github.com/vulkan-go/vulkan"
)

var GlobalPool Pool

func InitGlobalPool(device device.T) {
	GlobalPool = NewPool(device, []vk.DescriptorPoolSize{
		{
			Type:            vk.DescriptorTypeUniformBuffer,
			DescriptorCount: 1000,
		},
		{
			Type:            vk.DescriptorTypeStorageBuffer,
			DescriptorCount: 1000,
		},
		{
			Type:            vk.DescriptorTypeCombinedImageSampler,
			DescriptorCount: 80,
		},
		{
			Type:            vk.DescriptorTypeInputAttachment,
			DescriptorCount: 100,
		},
	})
}

func DestroyGlobalPool() {
	GlobalPool.Destroy()
	GlobalPool = nil
}
