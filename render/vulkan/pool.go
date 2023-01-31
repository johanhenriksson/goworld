package vulkan

import (
	vk "github.com/vulkan-go/vulkan"
)

var DefaultDescriptorPools = []vk.DescriptorPoolSize{
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
		DescriptorCount: 10000,
	},
	{
		Type:            vk.DescriptorTypeInputAttachment,
		DescriptorCount: 100,
	},
}
