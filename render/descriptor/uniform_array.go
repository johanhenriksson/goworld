package descriptor

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/samber/lo"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type UniformArray[K any] struct {
	Size   int
	Stages core1_0.ShaderStageFlags

	binding int
	buffer  *buffer.Array[K]
	set     Set
}

func (d *UniformArray[K]) Initialize(dev *device.Device) {
	if d.set == nil {
		panic("descriptor must be bound first")
	}
	d.buffer = buffer.NewArray[K](dev, buffer.Args{
		Key:    d.String(),
		Size:   d.Size,
		Usage:  core1_0.BufferUsageUniformBuffer,
		Memory: device.MemoryTypeShared,
	})
	d.write()
}

func (d *UniformArray[K]) String() string {
	var empty K
	kind := reflect.TypeOf(empty)
	return fmt.Sprintf("UniformArray[%s]:%d", kind.Name(), d.binding)
}

func (d *UniformArray[K]) Destroy() {
	if d.buffer != nil {
		d.buffer.Destroy()
		d.buffer = nil
	}
}

func (d *UniformArray[K]) Bind(set Set, binding int) {
	d.set = set
	d.binding = binding
}

func (d *UniformArray[K]) Set(index int, data K) {
	d.buffer.Set(index, data)
}

func (d *UniformArray[K]) SetRange(offset int, data []K) {
	d.buffer.SetRange(offset, data)
}

func (d *UniformArray[K]) write() {
	d.set.Write(core1_0.WriteDescriptorSet{
		DstBinding:      d.binding,
		DstArrayElement: 0,
		DescriptorType:  core1_0.DescriptorTypeUniformBuffer,
		BufferInfo: lo.Times(d.Size, func(i int) core1_0.DescriptorBufferInfo {
			return core1_0.DescriptorBufferInfo{
				Buffer: d.buffer.Ptr(),
				Offset: i * d.buffer.Stride(),
				Range:  d.buffer.Stride(),
			}
		}),
	})
}

func (d *UniformArray[K]) LayoutBinding(binding int) core1_0.DescriptorSetLayoutBinding {
	d.binding = binding
	return core1_0.DescriptorSetLayoutBinding{
		Binding:         binding,
		DescriptorType:  core1_0.DescriptorTypeUniformBuffer,
		DescriptorCount: d.Size,
		StageFlags:      core1_0.ShaderStageFlags(d.Stages),
	}
}

func (d *UniformArray[K]) BindingFlags() ext_descriptor_indexing.DescriptorBindingFlags { return 0 }
