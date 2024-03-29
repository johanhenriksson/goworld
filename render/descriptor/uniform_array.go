package descriptor

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/util"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type UniformArray[K any] struct {
	Size   int
	Stages core1_0.ShaderStageFlags

	binding int
	buffer  buffer.Array[K]
	set     Set
}

func (d *UniformArray[K]) Initialize(device device.T) {
	if d.set == nil {
		panic("descriptor must be bound first")
	}
	d.buffer = buffer.NewArray[K](device, buffer.Args{
		Key:    d.String(),
		Size:   d.Size,
		Usage:  core1_0.BufferUsageUniformBuffer,
		Memory: core1_0.MemoryPropertyDeviceLocal | core1_0.MemoryPropertyHostVisible,
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
		BufferInfo: util.Map(util.Range(0, d.Size, 1), func(i int) core1_0.DescriptorBufferInfo {
			return core1_0.DescriptorBufferInfo{
				Buffer: d.buffer.Ptr(),
				Offset: i * d.buffer.Element(),
				Range:  d.buffer.Element(),
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
