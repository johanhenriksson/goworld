package descriptor

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type UniformArray[K any] struct {
	Binding int
	Count   int
	Stages  vk.ShaderStageFlagBits

	element int
	buffer  buffer.T
	set     Set
}

func (d *UniformArray[K]) Initialize(device device.T) {
	if d.set == nil {
		panic("descriptor must be bound first")
	}

	var empty K
	t := reflect.TypeOf(empty)
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("UniformArray value must be a struct, was %s", t.Kind()))
	}

	d.element = (int(t.Size())/64 + 1) * 64
	d.buffer = buffer.NewUniform(device, d.element*d.Count)
	d.write()
}

func (d *UniformArray[K]) Destroy() {
	if d.buffer != nil {
		d.buffer.Destroy()
	}
}

func (d *UniformArray[K]) Bind(set Set) {
	d.set = set
}

func (d *UniformArray[K]) Set(index int, data K) {
	ptr := &data
	offset := index * d.element
	d.buffer.Write(ptr, offset)
}

func (d *UniformArray[K]) write() {
	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstBinding:      uint32(d.Binding),
		DstArrayElement: 0,
		DescriptorCount: uint32(d.Count),
		DescriptorType:  vk.DescriptorTypeUniformBuffer,
		PBufferInfo: util.Map(util.Range(0, d.Count, 1), func(i int) vk.DescriptorBufferInfo {
			return vk.DescriptorBufferInfo{
				Buffer: d.buffer.Ptr(),
				Offset: vk.DeviceSize(i * d.element),
				Range:  vk.DeviceSize(d.element),
			}
		}),
	})
}

func (d *UniformArray[K]) LayoutBinding() vk.DescriptorSetLayoutBinding {
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(d.Binding),
		DescriptorType:  vk.DescriptorTypeUniformBuffer,
		DescriptorCount: uint32(d.Count),
		StageFlags:      vk.ShaderStageFlags(d.Stages),
	}
}

func (d *UniformArray[K]) BindingFlags() vk.DescriptorBindingFlags { return 0 }
