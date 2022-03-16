package descriptor

import (
	"fmt"
	"reflect"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type Uniform[K any] struct {
	Binding int
	Stages  vk.ShaderStageFlagBits

	buffer buffer.T
	set    Set
}

func (d *Uniform[K]) Initialize(device device.T) {
	if d.set == nil {
		panic("descriptor must be bound first")
	}

	var empty K
	if err := util.ValidateAlignment(empty); err != nil {
		panic(fmt.Sprintf("illegal Uniform struct: %s", err))
	}

	t := reflect.TypeOf(empty)
	d.buffer = buffer.NewUniform(device, int(t.Size()))
	d.write()
}

func (d *Uniform[K]) String() string {
	var empty K
	kind := reflect.TypeOf(empty)
	return fmt.Sprintf("Uniform[%s]:%d", kind.Name(), d.Binding)
}

func (d *Uniform[K]) Destroy() {
	if d.buffer != nil {
		d.buffer.Destroy()
		d.buffer = nil
	}
}

func (d *Uniform[K]) Bind(set Set) {
	d.set = set
}

func (d *Uniform[K]) Set(data K) {
	ptr := &data
	d.buffer.Write(0, ptr)
}

func (d *Uniform[K]) write() {
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
		StageFlags:      vk.ShaderStageFlags(d.Stages),
	}
}

func (d *Uniform[K]) BindingFlags() vk.DescriptorBindingFlags { return 0 }
