package descriptor

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type UniformArray[K any] struct {
	Size   int
	Stages vk.ShaderStageFlagBits

	binding int
	buffer  buffer.Array[K]
	set     Set
}

func (d *UniformArray[K]) Initialize(device device.T) {
	if d.set == nil {
		panic("descriptor must be bound first")
	}
	d.buffer = buffer.NewArray[K](device, buffer.Args{
		Size:   d.Size,
		Usage:  vk.BufferUsageUniformBufferBit,
		Memory: vk.MemoryPropertyDeviceLocalBit | vk.MemoryPropertyHostVisibleBit,
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

func (d *UniformArray[K]) write() {
	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstBinding:      uint32(d.binding),
		DstArrayElement: 0,
		DescriptorCount: uint32(d.Size),
		DescriptorType:  vk.DescriptorTypeUniformBuffer,
		PBufferInfo: util.Map(util.Range(0, d.Size, 1), func(i int) vk.DescriptorBufferInfo {
			return vk.DescriptorBufferInfo{
				Buffer: d.buffer.Ptr(),
				Offset: vk.DeviceSize(i * d.buffer.Element()),
				Range:  vk.DeviceSize(d.buffer.Element()),
			}
		}),
	})
}

func (d *UniformArray[K]) LayoutBinding(binding int) vk.DescriptorSetLayoutBinding {
	d.binding = binding
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(binding),
		DescriptorType:  vk.DescriptorTypeUniformBuffer,
		DescriptorCount: uint32(d.Size),
		StageFlags:      vk.ShaderStageFlags(d.Stages),
	}
}

func (d *UniformArray[K]) BindingFlags() vk.DescriptorBindingFlags { return 0 }
