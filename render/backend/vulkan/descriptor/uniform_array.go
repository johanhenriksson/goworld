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
	Size    int
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
	if err := ValidateShaderStruct(empty); err != nil {
		panic(fmt.Sprintf("illegal UniformArray struct: %s", err))
	}

	alignment := int(device.GetLimits().MinUniformBufferOffsetAlignment)
	maxSize := int(device.GetLimits().MaxUniformBufferRange)

	t := reflect.TypeOf(empty)
	d.element = util.Align(int(t.Size()), alignment)
	size := d.element * d.Size
	if size > maxSize {
		panic(fmt.Sprintf("uniform buffer too large: %d, max size: %d", size, maxSize))
	}

	d.buffer = buffer.NewUniform(device, d.element*d.Size)
	d.write()
}

func (d *UniformArray[K]) String() string {
	var empty K
	kind := reflect.TypeOf(empty)
	return fmt.Sprintf("UniformArray[%s]:%d", kind.Name(), d.Binding)
}

func (d *UniformArray[K]) Destroy() {
	if d.buffer != nil {
		d.buffer.Destroy()
		d.buffer = nil
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
		DescriptorCount: uint32(d.Size),
		DescriptorType:  vk.DescriptorTypeUniformBuffer,
		PBufferInfo: util.Map(util.Range(0, d.Size, 1), func(i int) vk.DescriptorBufferInfo {
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
		DescriptorCount: uint32(d.Size),
		StageFlags:      vk.ShaderStageFlags(d.Stages),
	}
}

func (d *UniformArray[K]) BindingFlags() vk.DescriptorBindingFlags { return 0 }
