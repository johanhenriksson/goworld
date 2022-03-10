package descriptor

import (
	"reflect"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Storage[K any] struct {
	Binding int
	Stages  vk.ShaderStageFlags
	Size    int

	buffer  buffer.T
	set     Set
	element int
}

func (d *Storage[K]) Initialize(device device.T) {
	if d.set == nil {
		panic("descriptor must be bound first")
	}
	if d.Size == 0 {
		panic("storage descriptor size must be set")
	}

	var empty K
	t := reflect.TypeOf(empty)
	d.element = int(t.Size())

	d.buffer = buffer.NewStorage(device, d.Size*d.element)
	d.write()
}

func (d *Storage[K]) Destroy() {
	if d.buffer != nil {
		d.buffer.Destroy()
		d.buffer = nil
	}
}

func (d *Storage[K]) Bind(set Set) {
	d.set = set
}

func (d *Storage[K]) Set(index int, data K) {
	ptr := &data
	offset := index * d.element
	d.buffer.Write(ptr, offset)
}

func (d *Storage[K]) SetRange(data []K, offset int) {
	offset *= d.element
	d.buffer.Write(data, offset)
}

func (d *Storage[K]) LayoutBinding() vk.DescriptorSetLayoutBinding {
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(d.Binding),
		DescriptorType:  vk.DescriptorTypeStorageBuffer,
		DescriptorCount: 1,
		StageFlags:      d.Stages,
	}
}

func (d *Storage[K]) write() {
	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstSet:          d.set.Ptr(),
		DstBinding:      uint32(d.Binding),
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
