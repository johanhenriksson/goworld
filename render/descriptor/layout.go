package descriptor

import (
	"log"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/shader"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/core1_2"
)

type Map map[string]Descriptor

type SetLayout interface {
	device.Resource[core1_0.DescriptorSetLayout]
	VariableCount() int
}

type SetLayoutTyped[S Set] interface {
	SetLayout
	Instantiate(pool Pool) S
}

type layout[S Set] struct {
	device    device.T
	shader    shader.T
	ptr       core1_0.DescriptorSetLayout
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
	createFlags := core1_0.DescriptorSetLayoutCreateFlags(0)
	bindings := make([]core1_0.DescriptorSetLayoutBinding, 0, len(descriptors))
	bindFlags := make([]core1_2.DescriptorBindingFlags, 0, len(descriptors))
	for name, descriptor := range descriptors {
		index, exists := shader.Descriptor(name)
		if !exists {
			panic("unresolved descriptor")
		}
		bindings = append(bindings, descriptor.LayoutBinding(index))
		flags := descriptor.BindingFlags()
		bindFlags = append(bindFlags, flags)

		if flags&core1_2.DescriptorBindingUpdateAfterBind > 0 {
			createFlags |= core1_2.DescriptorSetLayoutCreateUpdateAfterBindPool
		}

		log.Printf("  %s -> %s\n", name, descriptor)
		if variable, ok := descriptor.(VariableDescriptor); ok {
			maxCount = variable.MaxCount()
			log.Println("descriptor", name, "is of variable length", maxCount)
		}
	}

	bindFlagsInfo := core1_2.DescriptorSetLayoutBindingFlagsCreateInfo{
		BindingFlags: bindFlags,
	}

	info := core1_0.DescriptorSetLayoutCreateInfo{
		Flags:       core1_0.DescriptorSetLayoutCreateFlags(createFlags),
		Bindings:    bindings,
		NextOptions: bindFlagsInfo.NextOptions,
	}

	ptr, _, err := device.Ptr().CreateDescriptorSetLayout(nil, info)
	if err != nil {
		panic(err)
	}

	return &layout[S]{
		device:   device,
		shader:   shader,
		ptr:      ptr,
		set:      set,
		maxCount: maxCount,
	}
}

func (d *layout[S]) Ptr() core1_0.DescriptorSetLayout {
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
		d.ptr.Destroy(nil)
		d.ptr = nil
	}
}
