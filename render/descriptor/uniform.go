package descriptor

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/buffer"
	"github.com/johanhenriksson/goworld/render/device"

	vk "github.com/vulkan-go/vulkan"
)

type Uniform[K any] struct {
	Stages vk.ShaderStageFlagBits

	binding int
	buffer  buffer.Item[K]
	set     Set
}

func (d *Uniform[K]) Initialize(device device.T) {
	if d.set == nil {
		panic("descriptor must be bound first")
	}
	d.buffer = buffer.NewItem[K](device, buffer.Args{
		Usage:  vk.BufferUsageUniformBufferBit,
		Memory: vk.MemoryPropertyDeviceLocalBit | vk.MemoryPropertyHostVisibleBit,
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
	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstBinding:      uint32(d.binding),
		DstArrayElement: 0,
		DescriptorCount: 1,
		DescriptorType:  vk.DescriptorTypeUniformBuffer,
		PBufferInfo: []vk.DescriptorBufferInfo{
			{
				Buffer: d.buffer.Ptr(),
				Offset: 0,
				Range:  vk.DeviceSize(vk.WholeSize),
			},
		},
	})
}

func (d *Uniform[K]) LayoutBinding(binding int) vk.DescriptorSetLayoutBinding {
	d.binding = binding
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(binding),
		DescriptorType:  vk.DescriptorTypeUniformBuffer,
		DescriptorCount: 1,
		StageFlags:      vk.ShaderStageFlags(d.Stages),
	}
}

func (d *Uniform[K]) BindingFlags() vk.DescriptorBindingFlags { return 0 }
