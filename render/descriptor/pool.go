package descriptor

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/common"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type Pool struct {
	ptr     core1_0.DescriptorPool
	device  *device.Device
	sizes   []core1_0.DescriptorPoolSize
	maxSets int

	allocatedSets   int
	allocatedCounts map[core1_0.DescriptorType]int
}

func NewPool(device *device.Device, sets int, sizes []core1_0.DescriptorPoolSize) *Pool {
	p := &Pool{
		device:          device,
		ptr:             nil,
		sizes:           sizes,
		maxSets:         sets,
		allocatedCounts: make(map[core1_0.DescriptorType]int),
	}
	p.Recreate()
	return p
}

func (p *Pool) Ptr() core1_0.DescriptorPool {
	return p.ptr
}

func (p *Pool) Recreate() {
	p.Destroy()

	info := core1_0.DescriptorPoolCreateInfo{
		Flags: core1_0.DescriptorPoolCreateFreeDescriptorSet |
			ext_descriptor_indexing.DescriptorPoolCreateUpdateAfterBind,
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

func (p *Pool) Destroy() {
	if p.ptr == nil {
		return
	}
	p.ptr.Destroy(nil)
	p.ptr = nil
}

func (p *Pool) Allocate(layout SetLayout) Set {
	sets := p.AllocateMany(layout, 1)
	return sets[0]
}

func (p *Pool) AllocateMany(layout SetLayout, count int) []Set {
	layouts := make([]core1_0.DescriptorSetLayout, count)
	for i := range layouts {
		layouts[i] = layout.Ptr()
	}
	info := core1_0.DescriptorSetAllocateInfo{
		DescriptorPool: p.ptr,
		SetLayouts:     layouts,
	}

	if layout.VariableCount() > 0 {
		variableCounts := make([]int, count)
		for i := range variableCounts {
			variableCounts[i] = layout.VariableCount()
		}
		variableInfo := ext_descriptor_indexing.DescriptorSetVariableDescriptorCountAllocateInfo{
			DescriptorCounts: variableCounts,
		}
		info.NextOptions = common.NextOptions{Next: variableInfo}
	}

	ptrs, r, err := p.device.Ptr().AllocateDescriptorSets(info)
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

	sets := make([]Set, count)
	for i, ptr := range ptrs {
		p.device.SetDebugObjectName(
			driver.VulkanHandle(ptr.Handle()),
			core1_0.ObjectTypeDescriptorSet,
			fmt.Sprintf("%s:%d", layout.Name(), i))

		p.allocatedSets++
		for kind, count := range layout.Counts() {
			current, _ := p.allocatedCounts[kind]
			p.allocatedCounts[kind] = current + count
		}

		sets[i] = &set{
			device: p.device,
			ptr:    ptr,
			layout: layout,
		}
	}
	return sets
}
