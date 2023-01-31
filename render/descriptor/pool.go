package descriptor

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/render/device"

	vk "github.com/vulkan-go/vulkan"
)

type Pool interface {
	device.Resource[vk.DescriptorPool]

	Allocate(layouts SetLayout) Set
	Recreate()
}

type pool struct {
	ptr     vk.DescriptorPool
	device  device.T
	sizes   []vk.DescriptorPoolSize
	maxSets int
}

func NewPool(device device.T, sizes []vk.DescriptorPoolSize) Pool {
	p := &pool{
		device:  device,
		ptr:     vk.NullDescriptorPool,
		sizes:   sizes,
		maxSets: 100,
	}
	p.Recreate()
	return p
}

func (p *pool) Ptr() vk.DescriptorPool {
	return p.ptr
}

func (p *pool) Recreate() {
	p.Destroy()

	info := vk.DescriptorPoolCreateInfo{
		SType:         vk.StructureTypeDescriptorPoolCreateInfo,
		Flags:         vk.DescriptorPoolCreateFlags(vk.DescriptorPoolCreateUpdateAfterBindBit),
		PPoolSizes:    p.sizes,
		PoolSizeCount: uint32(len(p.sizes)),
		MaxSets:       uint32(p.maxSets),
	}
	var ptr vk.DescriptorPool
	if ok := vk.CreateDescriptorPool(p.device.Ptr(), &info, nil, &ptr); ok != vk.Success {
		panic("failed to create descriptor pool")
	}
	p.ptr = ptr
}

func (p *pool) Destroy() {
	if p.ptr == vk.NullDescriptorPool {
		return
	}
	vk.DestroyDescriptorPool(p.device.Ptr(), p.ptr, nil)
	p.ptr = vk.NullDescriptorPool
}

func (p *pool) Allocate(layout SetLayout) Set {
	info := vk.DescriptorSetAllocateInfo{
		SType:              vk.StructureTypeDescriptorSetAllocateInfo,
		DescriptorPool:     p.ptr,
		DescriptorSetCount: 1,
		PSetLayouts:        []vk.DescriptorSetLayout{layout.Ptr()},
	}

	if layout.VariableCount() > 0 {
		variableInfo := vk.DescriptorSetVariableDescriptorCountAllocateInfo{
			SType:              vk.StructureTypeDescriptorSetVariableDescriptorCountAllocateInfo,
			DescriptorSetCount: 1,
			PDescriptorCounts:  []uint32{uint32(layout.VariableCount())},
		}
		info.PNext = unsafe.Pointer(&variableInfo)
	}

	var ptr vk.DescriptorSet
	r := vk.AllocateDescriptorSets(p.device.Ptr(), &info, &ptr)
	if r != vk.Success {
		if r == vk.ErrorOutOfPoolMemory {
			panic("failed to allocate descriptor set: out of pool memory")
		}
		panic("failed to allocate descriptor set")
	}

	return &set{
		device: p.device,
		ptr:    ptr,
		layout: layout,
	}
}
