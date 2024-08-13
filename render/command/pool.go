package command

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/samber/lo"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type Pool struct {
	ptr    core1_0.CommandPool
	device *device.Device
}

func NewPool(device *device.Device, flags core1_0.CommandPoolCreateFlags, queueFamilyIdx int) *Pool {
	ptr, _, err := device.Ptr().CreateCommandPool(nil, core1_0.CommandPoolCreateInfo{
		Flags:            flags,
		QueueFamilyIndex: queueFamilyIdx,
	})
	if err != nil {
		panic(err)
	}
	return &Pool{
		ptr:    ptr,
		device: device,
	}
}

func (p *Pool) Ptr() core1_0.CommandPool {
	return p.ptr
}

func (p *Pool) Destroy() {
	p.ptr.Destroy(nil)
	p.ptr = nil
}

func (p *Pool) Allocate(level core1_0.CommandBufferLevel) *Buffer {
	buffers := p.AllocateBuffers(level, 1)
	return buffers[0]
}

func (p *Pool) AllocateBuffers(level core1_0.CommandBufferLevel, count int) []*Buffer {
	ptrs, _, err := p.device.Ptr().AllocateCommandBuffers(core1_0.CommandBufferAllocateInfo{
		CommandPool:        p.ptr,
		Level:              level,
		CommandBufferCount: count,
	})
	if err != nil {
		panic(err)
	}

	return lo.Map(ptrs, func(ptr core1_0.CommandBuffer, _ int) *Buffer {
		return newBuffer(p.device, p.ptr, ptr)
	})
}
