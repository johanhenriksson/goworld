package command

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type Pool interface {
	device.Resource[vk.CommandPool]

	Allocate(level vk.CommandBufferLevel) Buffer
	AllocateBuffers(level vk.CommandBufferLevel, count int) []Buffer
}

type pool struct {
	ptr    vk.CommandPool
	device device.T
}

func NewPool(device device.T, flags vk.CommandPoolCreateFlags, queueFlags vk.QueueFlags) Pool {
	queueIdx := device.GetQueueFamilyIndex(queueFlags)
	info := vk.CommandPoolCreateInfo{
		SType:            vk.StructureTypeCommandPoolCreateInfo,
		Flags:            flags,
		QueueFamilyIndex: uint32(queueIdx),
	}

	var ptr vk.CommandPool
	vk.CreateCommandPool(device.Ptr(), &info, nil, &ptr)

	return &pool{
		ptr:    ptr,
		device: device,
	}
}

func (p *pool) Ptr() vk.CommandPool {
	return p.ptr
}

func (p *pool) Destroy() {
	vk.DestroyCommandPool(p.device.Ptr(), p.ptr, nil)
}

func (p *pool) Allocate(level vk.CommandBufferLevel) Buffer {
	buffers := p.AllocateBuffers(level, 1)
	return buffers[0]
}

func (p *pool) AllocateBuffers(level vk.CommandBufferLevel, count int) []Buffer {
	info := vk.CommandBufferAllocateInfo{
		SType:              vk.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        p.ptr,
		Level:              level,
		CommandBufferCount: uint32(count),
	}

	ptrs := make([]vk.CommandBuffer, count)
	vk.AllocateCommandBuffers(p.device.Ptr(), &info, ptrs)

	return util.Map(ptrs, func(ptr vk.CommandBuffer) Buffer {
		return newBuffer(p.device, p.ptr, ptr)
	})
}
