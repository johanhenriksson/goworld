package descriptor

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Map map[string]Descriptor

type Layout interface {
	device.Resource[vk.DescriptorSetLayout]
}

type TypedLayout[S Set] interface {
	Layout
	Allocate() S
	Descriptor(name string) Descriptor
}

type Args struct {
	Descriptors map[string]Descriptor
}

type Descriptor interface {
	LayoutBinding() vk.DescriptorSetLayoutBinding
	Write()
	Bind(Set)
}

type layout[T any] struct {
	device      device.T
	ptr         vk.DescriptorSetLayout
	pool        Pool
	set         T
	descriptors Map
}

func New[S Set](device device.T, pool Pool, set S) TypedLayout[S] {
	descriptors := ParseDescriptors(set)

	bindings := make([]vk.DescriptorSetLayoutBinding, 0, len(descriptors))
	for _, descriptor := range descriptors {
		bindings = append(bindings, descriptor.LayoutBinding())
	}

	log.Println("descriptor bindings:", bindings)

	info := vk.DescriptorSetLayoutCreateInfo{
		SType:        vk.StructureTypeDescriptorSetLayoutCreateInfo,
		BindingCount: uint32(len(bindings)),
		PBindings:    bindings,
	}

	var ptr vk.DescriptorSetLayout
	vk.CreateDescriptorSetLayout(device.Ptr(), &info, nil, &ptr)

	return &layout[S]{
		device: device,
		ptr:    ptr,
		set:    set,
		pool:   pool,
	}
}

func (d *layout[T]) Ptr() vk.DescriptorSetLayout {
	return d.ptr
}

func (d *layout[T]) Allocate() T {
	dset := d.pool.AllocateSet(d)
	copy := BindSet(d.set, dset)
	return copy
}

func (d *layout[T]) Descriptor(name string) Descriptor {
	return d.descriptors[name]
}

func (d *layout[T]) Destroy() {
	if d.ptr != nil {
		vk.DestroyDescriptorSetLayout(d.device.Ptr(), d.ptr, nil)
		d.ptr = nil
	}
}

func ParseDescriptors(set any) Map {
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
	descriptors := make(Map)
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

		name := strings.ToLower(fieldName)
		descriptors[name] = descriptor
	}

	if !hasSet {
		panic("must embed descriptor.Set")
	}

	return descriptors
}

func BindSet[T any](set T, bind Set) T {
	// dereference
	ptr := reflect.ValueOf(set)
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
		}
	}

	copy := copyPtr.Interface().(T)
	return copy
}
