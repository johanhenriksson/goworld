package descriptor

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Storage[K any] struct {
	Stages vk.ShaderStageFlagBits
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
		Size:   d.Size,
		Usage:  vk.BufferUsageStorageBufferBit,
		Memory: vk.MemoryPropertyDeviceLocalBit | vk.MemoryPropertyHostVisibleBit,
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

func (d *Storage[K]) LayoutBinding(binding int) vk.DescriptorSetLayoutBinding {
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(binding),
		DescriptorType:  vk.DescriptorTypeStorageBuffer,
		DescriptorCount: 1,
		StageFlags:      vk.ShaderStageFlags(d.Stages),
	}
}

func (d *Storage[K]) BindingFlags() vk.DescriptorBindingFlags { return 0 }

func (d *Storage[K]) write() {
	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstSet:          d.set.Ptr(),
		DstBinding:      uint32(d.binding),
		DstArrayElement: 0,
		DescriptorCount: 1,
		DescriptorType:  vk.DescriptorTypeStorageBuffer,
		PBufferInfo: []vk.DescriptorBufferInfo{
			{
				Buffer: d.buffer.Ptr(),
				Offset: 0,
				Range:  vk.DeviceSize(vk.WholeSize),
			},
		},
	})
}
