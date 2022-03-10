package descriptor

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type Pool interface {
	device.Resource[vk.DescriptorPool]

	AllocateSet(layouts SetLayout) Set
	AllocateSets(layouts []SetLayout) []Set
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

func (p *pool) AllocateSet(layout SetLayout) Set {
	return p.AllocateSets([]SetLayout{layout})[0]
}

func (p *pool) AllocateSets(layouts []SetLayout) []Set {
	info := vk.DescriptorSetAllocateInfo{
		SType:              vk.StructureTypeDescriptorSetAllocateInfo,
		DescriptorPool:     p.ptr,
		DescriptorSetCount: uint32(len(layouts)),
		PSetLayouts: util.Map(layouts, func(item SetLayout) vk.DescriptorSetLayout {
			return item.Ptr()
		}),
	}

	sets := make([]vk.DescriptorSet, len(layouts))
	vk.AllocateDescriptorSets(p.device.Ptr(), &info, &sets[0])

	return util.MapIdx(sets, func(ptr vk.DescriptorSet, i int) Set {
		return &set{
			device: p.device,
			ptr:    ptr,
			layout: layouts[i],
		}
	})
}
