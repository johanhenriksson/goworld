package descriptor

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type Storage[K comparable] struct {
	Stages core1_0.ShaderStageFlags
	Size   int

	binding int
	buffer  buffer.Array[K]
	set     Set
}

func (d *Storage[K]) Initialize(device device.T) {
	if d.set == nil {
		panic("descriptor must be bound first")
	}
	if d.Size == 0 {
		panic("storage descriptor size must be non-zero")
	}

	d.buffer = buffer.NewArray[K](device, buffer.Args{
		Key:    d.String(),
		Size:   d.Size,
		Usage:  core1_0.BufferUsageStorageBuffer,
		Memory: core1_0.MemoryPropertyDeviceLocal | core1_0.MemoryPropertyHostVisible,
	})
	d.write()
}

func (d *Storage[K]) String() string {
	var empty K
	kind := reflect.TypeOf(empty)
	return fmt.Sprintf("Storage[%s]:%d", kind.Name(), d.binding)
}

func (d *Storage[K]) Destroy() {
	if d.buffer != nil {
		d.buffer.Destroy()
		d.buffer = nil
	}
}

func (d *Storage[K]) Bind(set Set, binding int) {
	d.set = set
	d.binding = binding
}

func (d *Storage[K]) Set(index int, data K) {
	d.buffer.Set(index, data)
}

func (d *Storage[K]) SetRange(offset int, data []K) {
	d.buffer.SetRange(offset, data)
}

func (d *Storage[K]) LayoutBinding(binding int) core1_0.DescriptorSetLayoutBinding {
	d.binding = binding
	return core1_0.DescriptorSetLayoutBinding{
		Binding:         binding,
		DescriptorType:  core1_0.DescriptorTypeStorageBuffer,
		DescriptorCount: 1,
		StageFlags:      core1_0.ShaderStageFlags(d.Stages),
	}
}

func (d *Storage[K]) BindingFlags() ext_descriptor_indexing.DescriptorBindingFlags { return 0 }

func (d *Storage[K]) write() {
	d.set.Write(core1_0.WriteDescriptorSet{
		DstSet:          d.set.Ptr(),
		DstBinding:      d.binding,
		DstArrayElement: 0,
		DescriptorType:  core1_0.DescriptorTypeStorageBuffer,
		BufferInfo: []core1_0.DescriptorBufferInfo{
			{
				Buffer: d.buffer.Ptr(),
				Offset: 0,
				Range:  d.buffer.Size(),
			},
		},
	})
}
