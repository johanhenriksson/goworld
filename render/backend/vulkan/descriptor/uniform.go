package descriptor

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Uniform[K any] struct {
	Binding int
	Stages  vk.ShaderStageFlags

	buffer buffer.T
	set    Set
}

func (d *Uniform[K]) Initialize(device device.T) {
	if d.set == nil {
		panic("descriptor must be bound first")
	}

	var empty K
	t := reflect.TypeOf(empty)
	d.buffer = buffer.NewUniform(device, int(t.Size()))
	d.Write()
	fmt.Println("initialize uniform")
}

func (d *Uniform[K]) Destroy() {
	if d.buffer != nil {
		d.buffer.Destroy()
	}
}

func (d *Uniform[K]) Bind(set Set) {
	d.set = set
}

func (d *Uniform[K]) Set(data K) {
	ptr := &data
	d.buffer.Write(ptr, 0)
}

func (d *Uniform[K]) Write() {
	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstBinding:      uint32(d.Binding),
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

func (d *Uniform[K]) LayoutBinding() vk.DescriptorSetLayoutBinding {
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(d.Binding),
		DescriptorType:  vk.DescriptorTypeUniformBuffer,
		DescriptorCount: 1,
		StageFlags:      d.Stages,
	}
}
