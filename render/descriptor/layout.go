package descriptor

import (
	"log"
	"unsafe"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/shader"

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
	shader    shader.T
	ptr       vk.DescriptorSetLayout
	set       S
	allocated []Descriptor
	maxCount  int
}

func New[S Set](device device.T, set S, shader shader.T) SetLayoutTyped[S] {
	descriptors, err := ParseDescriptorStruct(set)
	if err != nil {
		panic(err)
	}

	log.Println("descriptor set")
	maxCount := 0
	createFlags := vk.DescriptorSetLayoutCreateFlagBits(0)
	bindings := make([]vk.DescriptorSetLayoutBinding, 0, len(descriptors))
	bindFlags := make([]vk.DescriptorBindingFlags, 0, len(descriptors))
	for name, descriptor := range descriptors {
		index, exists := shader.Descriptor(name)
		if !exists {
			panic("unresolved descriptor")
		}
		bindings = append(bindings, descriptor.LayoutBinding(index))
		flags := descriptor.BindingFlags()
		bindFlags = append(bindFlags, flags)

		if flags&vk.DescriptorBindingFlags(vk.DescriptorBindingUpdateAfterBindBit) > 0 {
			createFlags |= vk.DescriptorSetLayoutCreateUpdateAfterBindPoolBit
		}

		log.Printf("  %s -> %s\n", name, descriptor)
		if variable, ok := descriptor.(VariableDescriptor); ok {
			maxCount = variable.MaxCount()
			log.Println("descriptor", name, "is of variable length", maxCount)
		}
	}

	bindFlagsInfo := vk.DescriptorSetLayoutBindingFlagsCreateInfo{
		SType:         vk.StructureTypeDescriptorSetLayoutBindingFlagsCreateInfo,
		BindingCount:  uint32(len(bindFlags)),
		PBindingFlags: bindFlags,
	}

	info := vk.DescriptorSetLayoutCreateInfo{
		SType: vk.StructureTypeDescriptorSetLayoutCreateInfo,
		PNext: unsafe.Pointer(&bindFlagsInfo),

		Flags:        vk.DescriptorSetLayoutCreateFlags(createFlags),
		BindingCount: uint32(len(bindings)),
		PBindings:    bindings,
	}

	var ptr vk.DescriptorSetLayout
	vk.CreateDescriptorSetLayout(device.Ptr(), &info, nil, &ptr)

	return &layout[S]{
		device:   device,
		shader:   shader,
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
	copy, descriptors := CopyDescriptorStruct(d.set, set, d.shader)
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
