package descriptor

import (
	"fmt"
	"log"
	"reflect"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/buffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

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
	t := reflect.TypeOf(empty)
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("Uniform value must be a struct, was %s", t.Kind()))
	}

	log.Println("uniform of type", t.Name())
	expectedOffset := 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		log.Println("  field", field.Name, "offset:", field.Offset, "size:", field.Type.Size())
		if field.Offset != uintptr(expectedOffset) {
			panic("struct layout causes alignment issues")
		}
		expectedOffset = int(field.Offset + field.Type.Size())
	}

	d.buffer = buffer.NewUniform(device, int(t.Size()))
	d.write()
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
