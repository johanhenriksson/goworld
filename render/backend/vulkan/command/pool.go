package command

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type Pool interface {
	device.Resource
	Ptr() vk.CommandPool
}

type pool struct {
	ptr    vk.CommandPool
	device device.T
}

func New(device device.T, queueFamily int) Pool {
	info := vk.CommandPoolCreateInfo{
		SType:            vk.StructureTypeCommandPoolCreateInfo,
		QueueFamilyIndex: uint32(queueFamily),
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

func (p *pool) AllocateBuffers(count int) []Buffer {
	info := vk.CommandBufferAllocateInfo{
		SType: vk.StructureTypeCommandBufferAllocateInfo,
		Level: vk.CommandBufferLevelPrimary,
	}

	ptrs := make([]vk.CommandBuffer, count)
	vk.AllocateCommandBuffers(p.device.Ptr(), &info, ptrs)

	return util.Map(ptrs, func(i int, ptr vk.CommandBuffer) Buffer {
		return newBuffer(p.device, ptr)
	})
}
