package descriptor

import (
	"log"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/shader"

	"github.com/vkngwrapper/core/v2/common"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type Map map[string]Descriptor

type SetLayout interface {
	device.Resource[core1_0.DescriptorSetLayout]
	Name() string
	Counts() map[core1_0.DescriptorType]int
	VariableCount() int
}

type Layout[S Set] struct {
	device    *device.Device
	shader    *shader.Shader
	ptr       core1_0.DescriptorSetLayout
	set       S
	allocated []Descriptor
	maxCount  int
	counts    map[core1_0.DescriptorType]int
}

func New[S Set](device *device.Device, set S, shader *shader.Shader) *Layout[S] {
	descriptors, err := ParseDescriptorStruct(set)
	if err != nil {
		panic(err)
	}

	log.Println("descriptor set")
	maxCount := 0
	createFlags := core1_0.DescriptorSetLayoutCreateFlags(0)
	bindings := make([]core1_0.DescriptorSetLayoutBinding, 0, len(descriptors))
	bindFlags := make([]ext_descriptor_indexing.DescriptorBindingFlags, 0, len(descriptors))
	counts := make(map[core1_0.DescriptorType]int)
	for name, descriptor := range descriptors {
		index, exists := shader.Descriptor(name)
		if !exists {
			panic("unresolved descriptor")
		}
		binding := descriptor.LayoutBinding(index)
		bindings = append(bindings, binding)
		flags := descriptor.BindingFlags()
		bindFlags = append(bindFlags, flags)

		// always allow updating while pending
		flags = flags | ext_descriptor_indexing.DescriptorBindingUpdateUnusedWhilePending

		if flags&ext_descriptor_indexing.DescriptorBindingUpdateAfterBind == ext_descriptor_indexing.DescriptorBindingUpdateAfterBind {
			createFlags |= ext_descriptor_indexing.DescriptorSetLayoutCreateUpdateAfterBindPool
		}

		if variable, ok := descriptor.(VariableDescriptor); ok {
			maxCount = variable.MaxCount()
			log.Printf("  %s -> %s x0-%d\n", name, descriptor, maxCount)
			counts[binding.DescriptorType] = maxCount
		} else {
			log.Printf("  %s -> %s x%d\n", name, descriptor, binding.DescriptorCount)
			counts[binding.DescriptorType] = binding.DescriptorCount
		}
	}

	bindFlagsInfo := ext_descriptor_indexing.DescriptorSetLayoutBindingFlagsCreateInfo{
		BindingFlags: bindFlags,
	}

	info := core1_0.DescriptorSetLayoutCreateInfo{
		Flags:       createFlags,
		Bindings:    bindings,
		NextOptions: common.NextOptions{Next: bindFlagsInfo},
	}

	ptr, _, err := device.Ptr().CreateDescriptorSetLayout(nil, info)
	if err != nil {
		panic(err)
	}

	device.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()), core1_0.ObjectTypeDescriptorSetLayout, shader.Name())

	return &Layout[S]{
		device:   device,
		shader:   shader,
		ptr:      ptr,
		set:      set,
		maxCount: maxCount,
		counts:   counts,
	}
}

func (d *Layout[S]) Name() string {
	return d.shader.Name()
}

func (d *Layout[S]) Ptr() core1_0.DescriptorSetLayout {
	return d.ptr
}

func (d *Layout[S]) Counts() map[core1_0.DescriptorType]int {
	return d.counts
}

func (d *Layout[S]) VariableCount() int {
	return d.maxCount
}

func (d *Layout[S]) Instantiate(pool *Pool) S {
	set := pool.Allocate(d)
	copy, descriptors := CopyDescriptorStruct(d.set, set, d.shader)
	for _, descriptor := range descriptors {
		descriptor.Initialize(d.device)
		d.allocated = append(d.allocated, descriptor)
	}
	return copy
}

func (d *Layout[S]) Destroy() {
	// todo: allocated sets should probably clean up themselves
	for _, desc := range d.allocated {
		desc.Destroy()
	}
	if d.ptr != nil {
		d.ptr.Destroy(nil)
		d.ptr = nil
	}
}
