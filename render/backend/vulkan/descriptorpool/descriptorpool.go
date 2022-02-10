package descriptorpool

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource
	Ptr() vk.DescriptorPool
}

type pool struct {
	ptr    vk.DescriptorPool
	device device.T
}

func New(device device.T, sizes []vk.DescriptorPoolSize) T {
	info := vk.DescriptorPoolCreateInfo{
		SType:         vk.StructureTypeDescriptorPoolCreateInfo,
		PPoolSizes:    sizes,
		PoolSizeCount: uint32(len(sizes)),
	}

	var ptr vk.DescriptorPool
	vk.CreateDescriptorPool(device.Ptr(), &info, nil, &ptr)

	return &pool{
		device: device,
		ptr:    ptr,
	}
}

func (p *pool) Ptr() vk.DescriptorPool {
	return p.ptr
}

func (p *pool) Destroy() {
	vk.DestroyDescriptorPool(p.device.Ptr(), p.ptr, nil)
	p.ptr = nil
}
