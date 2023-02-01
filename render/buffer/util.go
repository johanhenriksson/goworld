package buffer

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

func GetBufferLimits(device device.T, usage core1_0.BufferUsageFlags) (align, max int) {
	limits := device.GetLimits()
	if usage&core1_0.BufferUsageUniformBuffer > 0 {
		return int(limits.MinUniformBufferOffsetAlignment), int(limits.MaxUniformBufferRange)
	}
	if usage&core1_0.BufferUsageStorageBuffer > 0 {
		return int(limits.MinStorageBufferOffsetAlignment), int(limits.MaxStorageBufferRange)
	}
	panic("unknown buffer usage type")
}
