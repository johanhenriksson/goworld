package descriptor

import (
	"unsafe"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Pool interface {
	device.Resource[vk.DescriptorPool]

	Allocate(layouts SetLayout) Set
}

type pool struct {
	ptr    vk.DescriptorPool
	device device.T
}

func NewPool(device device.T, sizes []vk.DescriptorPoolSize) Pool {
	info := vk.DescriptorPoolCreateInfo{
		SType:         vk.StructureTypeDescriptorPoolCreateInfo,
		Flags:         vk.DescriptorPoolCreateFlags(vk.DescriptorPoolCreateUpdateAfterBindBit),
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
		panic("failed to allocate descriptor set")
	}

	return &set{
		device: p.device,
		ptr:    ptr,
		layout: layout,
	}
}
