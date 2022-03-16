package buffer

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

func GetBufferLimits(device device.T, usage vk.BufferUsageFlagBits) (align, max int) {
	limits := device.GetLimits()
	if usage&vk.BufferUsageUniformBufferBit > 0 {
		return int(limits.MinUniformBufferOffsetAlignment), int(limits.MaxUniformBufferRange)
	}
	if usage&vk.BufferUsageStorageBufferBit > 0 {
		return int(limits.MinStorageBufferOffsetAlignment), int(limits.MaxStorageBufferRange)
	}
	panic("unknown buffer usage type")
}
