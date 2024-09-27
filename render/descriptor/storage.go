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
	Buffer *buffer.Array[K]

	binding  int
	set      Set
	ownedbuf bool
}

var _ Descriptor = (*Storage[any])(nil)

func (d *Storage[K]) Initialize(dev *device.Device, set Set, binding int) {
	d.set = set
	d.binding = binding

	if d.Buffer == nil {
		if d.Size == 0 {
			panic("storage descriptor size must be non-zero")
		}
		d.Buffer = buffer.NewArray[K](dev, buffer.Args{
			Key:    d.String(),
			Size:   d.Size,
			Usage:  core1_0.BufferUsageStorageBuffer,
			Memory: device.MemoryTypeShared,
		})
		d.ownedbuf = true
	} else {
		if d.Size > 0 && d.Size != d.Buffer.Size() {
			panic("storage descriptor size mismatch")
		}
		d.Size = d.Buffer.Size()
	}
	d.write()
}

func (d *Storage[K]) String() string {
	var empty K
	kind := reflect.TypeOf(empty)
	return fmt.Sprintf("Storage[%s]:%d", kind.Name(), d.binding)
}

func (d *Storage[K]) Destroy() {
	if d.Buffer != nil && d.ownedbuf {
		d.Buffer.Destroy()
		d.Buffer = nil
		d.ownedbuf = false
	}
}

func (d *Storage[K]) Set(index int, data K) {
	d.Buffer.Set(index, data)
}

func (d *Storage[K]) SetRange(offset int, data []K) {
	d.Buffer.SetRange(offset, data)
}

func (d *Storage[K]) LayoutBinding(binding int) core1_0.DescriptorSetLayoutBinding {
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
				Buffer: d.Buffer.Ptr(),
				Offset: 0,
				Range:  d.Buffer.Size(),
			},
		},
	})
}
