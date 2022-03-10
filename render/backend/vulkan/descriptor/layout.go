package descriptor

import (
	"fmt"
	"log"
	"reflect"
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
	Allocate() S
}

type layout[S any] struct {
	device      device.T
	ptr         vk.DescriptorSetLayout
	pool        Pool
	set         S
	allocated   []Descriptor
	descriptors []Descriptor
}

func New[S Set](device device.T, pool Pool, set S) SetLayoutTyped[S] {
	descriptors := ParseDescriptors(set)

	bindings := util.Map(descriptors, func(desc Descriptor) vk.DescriptorSetLayoutBinding {
		return desc.LayoutBinding()
	})

	bindFlags := util.Map(descriptors, func(desc Descriptor) vk.DescriptorBindingFlags { return desc.BindingFlags() })

	log.Println("descriptor bindings:", bindings)

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
		device:      device,
		ptr:         ptr,
		set:         set,
		pool:        pool,
		descriptors: descriptors,
	}
}

func (d *layout[S]) Ptr() vk.DescriptorSetLayout {
	return d.ptr
}

func (d *layout[S]) VariableCount() int {
	lastDescriptor := d.descriptors[len(d.descriptors)-1]
	if variable, ok := lastDescriptor.(VariableDescriptor); ok {
		return variable.MaxCount()
	}
	return 0
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

func ParseDescriptors[S Set](set S) []Descriptor {
	ptr := reflect.ValueOf(set)
	if ptr.Type().Kind() != reflect.Pointer {
		panic("set is not a pointer to struct")
	}

	// dereference pointer
	value := ptr.Elem()

	if value.Kind() != reflect.Struct {
		panic(fmt.Sprintf("set d is not a pointer to struct, was %s", value.Kind()))
	}

	hasSet := false
	descriptors := make([]Descriptor, 0, value.NumField())
	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name

		if fieldName == "Set" {
			hasSet = true
			continue
		}

		field := value.Field(i)
		descriptor, ok := field.Interface().(Descriptor)
		if !ok {
			panic(fmt.Sprintf("%s is not a Descriptor value\n", fieldName))
		}

		descriptors = append(descriptors, descriptor)
	}

	if !hasSet {
		panic("must embed descriptor.Set")
	}

	return descriptors
}

func (d *layout[S]) Allocate() S {
	bind := d.pool.Allocate(d)

	// dereference
	ptr := reflect.ValueOf(d.set)
	value := ptr.Elem()
	structName := value.Type().Name()

	copyPtr := reflect.New(value.Type())

	for i := 0; i < value.NumField(); i++ {
		fieldName := value.Type().Field(i).Name

		if fieldName == "Set" {
			copyPtr.Elem().Field(i).Set(reflect.ValueOf(bind))
		} else {
			field := value.Field(i)
			if field.Kind() != reflect.Pointer {
				panic(fmt.Sprintf("descriptor %s.%s is not a pointer, was %s", structName, fieldName, field.Kind()))
			}
			if field.IsZero() {
				panic(fmt.Sprintf("descriptor %s.%s is unset", structName, fieldName))
			}

			descValue := field.Elem()

			descCopy := reflect.New(descValue.Type())
			descCopy.Elem().Set(descValue)
			copyPtr.Elem().Field(i).Set(descCopy)

			desc := descCopy.Interface().(Descriptor)
			desc.Bind(bind)
			desc.Initialize(d.device)
			d.allocated = append(d.allocated, desc)
		}
	}

	copy := copyPtr.Interface().(S)
	return copy
}
