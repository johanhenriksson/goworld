package descriptor

import (
	"log"
	"unsafe"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type Map map[string]Descriptor

type SetLayout interface {
	device.Resource[vk.DescriptorSetLayout]
	VariableCount() int
}

type SetLayoutTyped[S Set] interface {
	SetLayout
	Instantiate(pool Pool) S
}

type layout[S Set] struct {
	device    device.T
	ptr       vk.DescriptorSetLayout
	set       S
	allocated []Descriptor
	maxCount  int
}

func New[S Set](device device.T, set S) SetLayoutTyped[S] {
	descriptors, err := ParseDescriptorStruct(set)
	if err != nil {
		panic(err)
	}

	log.Println("descriptors:", descriptors)

	maxCount := 0
	lastDescriptor := descriptors[len(descriptors)-1]
	if variable, ok := lastDescriptor.(VariableDescriptor); ok {
		log.Println("last descriptor is of variable length")
		maxCount = variable.MaxCount()
	}

	bindings := util.Map(descriptors, func(desc Descriptor) vk.DescriptorSetLayoutBinding {
		return desc.LayoutBinding()
	})

	bindFlags := util.Map(descriptors, func(desc Descriptor) vk.DescriptorBindingFlags { return desc.BindingFlags() })

	bindFlagsInfo := vk.DescriptorSetLayoutBindingFlagsCreateInfo{
		SType:         vk.StructureTypeDescriptorSetLayoutBindingFlagsCreateInfo,
		BindingCount:  uint32(len(bindFlags)),
		PBindingFlags: bindFlags,
	}

	info := vk.DescriptorSetLayoutCreateInfo{
		SType: vk.StructureTypeDescriptorSetLayoutCreateInfo,
		PNext: unsafe.Pointer(&bindFlagsInfo),

		BindingCount: uint32(len(bindings)),
		PBindings:    bindings,
	}

	var ptr vk.DescriptorSetLayout
	vk.CreateDescriptorSetLayout(device.Ptr(), &info, nil, &ptr)

	return &layout[S]{
		device:   device,
		ptr:      ptr,
		set:      set,
		maxCount: maxCount,
	}
}

func (d *layout[S]) Ptr() vk.DescriptorSetLayout {
	return d.ptr
}

func (d *layout[S]) VariableCount() int {
	return d.maxCount
}

func (d *layout[S]) Instantiate(pool Pool) S {
	set := pool.Allocate(d)
	copy, descriptors := CopyDescriptorStruct(d.set, set)
	for _, descriptor := range descriptors {
		descriptor.Initialize(d.device)
		d.allocated = append(d.allocated, descriptor)
	}
	return copy
}

func (d *layout[S]) Destroy() {
	for _, desc := range d.allocated {
		desc.Destroy()
	}
	if d.ptr != nil {
		vk.DestroyDescriptorSetLayout(d.device.Ptr(), d.ptr, nil)
		d.ptr = nil
	}
}
