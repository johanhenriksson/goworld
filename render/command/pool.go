package command

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/util"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type Pool interface {
	device.Resource[core1_0.CommandPool]

	Allocate(level core1_0.CommandBufferLevel) Buffer
	AllocateBuffers(level core1_0.CommandBufferLevel, count int) []Buffer
}

type pool struct {
	ptr    core1_0.CommandPool
	device device.T
}

func NewPool(device device.T, flags core1_0.CommandPoolCreateFlags, queueFamilyIdx int) Pool {
	ptr, _, err := device.Ptr().CreateCommandPool(nil, core1_0.CommandPoolCreateInfo{
		Flags:            flags,
		QueueFamilyIndex: queueFamilyIdx,
	})
	if err != nil {
		panic(err)
	}
	return &pool{
		ptr:    ptr,
		device: device,
	}
}

func (p *pool) Ptr() core1_0.CommandPool {
	return p.ptr
}

func (p *pool) Destroy() {
	p.ptr.Destroy(nil)
	p.ptr = nil
}

func (p *pool) Allocate(level core1_0.CommandBufferLevel) Buffer {
	buffers := p.AllocateBuffers(level, 1)
	return buffers[0]
}

func (p *pool) AllocateBuffers(level core1_0.CommandBufferLevel, count int) []Buffer {
	ptrs, _, err := p.device.Ptr().AllocateCommandBuffers(core1_0.CommandBufferAllocateInfo{
		CommandPool:        p.ptr,
		Level:              level,
		CommandBufferCount: count,
	})
	if err != nil {
		panic(err)
	}

	return util.Map(ptrs, func(ptr core1_0.CommandBuffer) Buffer {
		return newBuffer(p.device, p.ptr, ptr)
	})
}
