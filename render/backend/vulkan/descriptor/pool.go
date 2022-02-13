package descriptor

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type Pool interface {
	device.Resource[vk.DescriptorPool]

	AllocateSets(layouts []T) []Set
}

type pool struct {
	ptr    vk.DescriptorPool
	device device.T
}

func NewPool(device device.T, sizes []vk.DescriptorPoolSize) Pool {
	info := vk.DescriptorPoolCreateInfo{
		SType:         vk.StructureTypeDescriptorPoolCreateInfo,
		PPoolSizes:    sizes,
		PoolSizeCount: uint32(len(sizes)),
		MaxSets:       100,
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

func (p *pool) AllocateSets(layouts []T) []Set {
	info := vk.DescriptorSetAllocateInfo{
		SType:              vk.StructureTypeDescriptorSetAllocateInfo,
		DescriptorPool:     p.ptr,
		DescriptorSetCount: uint32(len(layouts)),
		PSetLayouts: util.Map(layouts, func(i int, item T) vk.DescriptorSetLayout {
			return item.Ptr()
		}),
	}

	sets := make([]vk.DescriptorSet, len(layouts))
	vk.AllocateDescriptorSets(p.device.Ptr(), &info, &sets[0])

	return util.Map(sets, func(i int, ptr vk.DescriptorSet) Set {
		return &set{
			ptr: ptr,
		}
	})
}
