package engine

import "github.com/vkngwrapper/core/v2/core1_0"

var DefaultDescriptorPools = []core1_0.DescriptorPoolSize{
	{
		Type:            core1_0.DescriptorTypeUniformBuffer,
		DescriptorCount: 10000,
	},
	{
		Type:            core1_0.DescriptorTypeStorageBuffer,
		DescriptorCount: 10000,
	},
	{
		Type:            core1_0.DescriptorTypeCombinedImageSampler,
		DescriptorCount: 100000,
	},
	{
		Type:            core1_0.DescriptorTypeInputAttachment,
		DescriptorCount: 1000,
	},
}
