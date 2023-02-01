package descriptor

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type Uniform[K any] struct {
	Stages core1_0.ShaderStageFlags

	binding int
	buffer  buffer.Item[K]
	set     Set
}

func (d *Uniform[K]) Initialize(device device.T) {
	if d.set == nil {
		panic("descriptor must be bound first")
	}
	d.buffer = buffer.NewItem[K](device, buffer.Args{
		Usage:  core1_0.BufferUsageUniformBuffer,
		Memory: core1_0.MemoryPropertyDeviceLocal | core1_0.MemoryPropertyHostVisible,
	})
	d.write()
}

func (d *Uniform[K]) String() string {
	var empty K
	kind := reflect.TypeOf(empty)
	return fmt.Sprintf("Uniform[%s]:%d", kind.Name(), d.binding)
}

func (d *Uniform[K]) Destroy() {
	if d.buffer != nil {
		d.buffer.Destroy()
		d.buffer = nil
	}
}

func (d *Uniform[K]) Bind(set Set, binding int) {
	d.set = set
	d.binding = binding
}

func (d *Uniform[K]) Set(data K) {
	d.buffer.Set(data)
}

func (d *Uniform[K]) write() {
	d.set.Write(core1_0.WriteDescriptorSet{
		DstBinding:      d.binding,
		DstArrayElement: 0,
		DescriptorType:  core1_0.DescriptorTypeUniformBuffer,
		BufferInfo: []core1_0.DescriptorBufferInfo{
			{
				Buffer: d.buffer.Ptr(),
				Offset: 0,
				Range:  d.buffer.Size(),
			},
		},
	})
}

func (d *Uniform[K]) LayoutBinding(binding int) core1_0.DescriptorSetLayoutBinding {
	d.binding = binding
	return core1_0.DescriptorSetLayoutBinding{
		Binding:         binding,
		DescriptorType:  core1_0.DescriptorTypeUniformBuffer,
		DescriptorCount: 1,
		StageFlags:      core1_0.ShaderStageFlags(d.Stages),
	}
}

func (d *Uniform[K]) BindingFlags() ext_descriptor_indexing.DescriptorBindingFlags { return 0 }
