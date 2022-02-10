package commandpool

import (
	"github.com/johanhenriksson/goworld/render/backend/utils"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/commandbuffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource
	Ptr() vk.CommandPool
}

type pool struct {
	ptr    vk.CommandPool
	device device.T
}

func New(device device.T, queueFamily int) T {
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

func (p *pool) AllocateBuffers(count int) []commandbuffer.T {
	info := vk.CommandBufferAllocateInfo{
		SType: vk.StructureTypeCommandBufferAllocateInfo,
		Level: vk.CommandBufferLevelPrimary,
	}

	ptrs := make([]vk.CommandBuffer, count)
	vk.AllocateCommandBuffers(p.device.Ptr(), &info, ptrs)

	return utils.Map(ptrs, func(ptr vk.CommandBuffer) commandbuffer.T {
		return commandbuffer.New(p.device, ptr)
	})
}
