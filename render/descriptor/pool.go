package descriptor

import (
	"log"

	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/common"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type Pool interface {
	device.Resource[core1_0.DescriptorPool]

	Allocate(layouts SetLayout) Set
	Recreate()
}

type pool struct {
	ptr     core1_0.DescriptorPool
	device  device.T
	sizes   []core1_0.DescriptorPoolSize
	maxSets int

	allocatedSets   int
	allocatedCounts map[core1_0.DescriptorType]int
}

func NewPool(device device.T, sets int, sizes []core1_0.DescriptorPoolSize) Pool {
	p := &pool{
		device:          device,
		ptr:             nil,
		sizes:           sizes,
		maxSets:         sets,
		allocatedCounts: make(map[core1_0.DescriptorType]int),
	}
	p.Recreate()
	return p
}

func (p *pool) Ptr() core1_0.DescriptorPool {
	return p.ptr
}

func (p *pool) Recreate() {
	p.Destroy()

	info := core1_0.DescriptorPoolCreateInfo{
		Flags:     ext_descriptor_indexing.DescriptorPoolCreateUpdateAfterBind,
		PoolSizes: p.sizes,
		MaxSets:   p.maxSets,
	}
	ptr, result, err := p.device.Ptr().CreateDescriptorPool(nil, info)
	if err != nil {
		panic(err)
	}
	if result != core1_0.VKSuccess {
		panic("failed to create descriptor pooll")
	}
	p.ptr = ptr
}

func (p *pool) Destroy() {
	if p.ptr == nil {
		return
	}
	p.ptr.Destroy(nil)
	p.ptr = nil
}

func (p *pool) Allocate(layout SetLayout) Set {
	info := core1_0.DescriptorSetAllocateInfo{
		DescriptorPool: p.ptr,
		SetLayouts:     []core1_0.DescriptorSetLayout{layout.Ptr()},
	}

	if layout.VariableCount() > 0 {
		variableInfo := ext_descriptor_indexing.DescriptorSetVariableDescriptorCountAllocateInfo{
			DescriptorCounts: []int{layout.VariableCount()},
		}
		info.NextOptions = common.NextOptions{Next: variableInfo}
	}

	ptr, r, err := p.device.Ptr().AllocateDescriptorSets(info)
	if err != nil {
		log.Println("allocated sets:", p.allocatedSets, "/", p.maxSets)
		log.Println("allocated counts:", p.allocatedCounts)
		panic(err)
	}
	if r != core1_0.VKSuccess {
		if r == core1_0.VKErrorOutOfDeviceMemory {
			panic("failed to allocate descriptor set: out of pool memory")
		}
		panic("failed to allocate descriptor set")
	}

	p.device.SetDebugObjectName(
		driver.VulkanHandle(ptr[0].Handle()),
		core1_0.ObjectTypeDescriptorSet,
		layout.Name())

	p.allocatedSets++
	for kind, count := range layout.Counts() {
		current, _ := p.allocatedCounts[kind]
		p.allocatedCounts[kind] = current + count
	}

	return &set{
		device: p.device,
		ptr:    ptr[0],
		layout: layout,
	}
}
