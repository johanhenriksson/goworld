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
	Buffer *buffer.Item[K]

	binding  int
	set      Set
	ownedbuf bool
}

var _ Descriptor = (*Uniform[any])(nil)

func (d *Uniform[K]) Initialize(dev *device.Device, set Set, binding int) {
	d.set = set
	d.binding = binding

	if d.Buffer == nil {
		d.Buffer = buffer.NewItem[K](dev, buffer.Args{
			Usage:  core1_0.BufferUsageUniformBuffer,
			Memory: device.MemoryTypeShared,
		})
		d.ownedbuf = true
	}
	d.write()
}

func (d *Uniform[K]) String() string {
	var empty K
	kind := reflect.TypeOf(empty)
	return fmt.Sprintf("Uniform[%s]:%d", kind.Name(), d.binding)
}

func (d *Uniform[K]) Destroy() {
	if d.Buffer != nil && d.ownedbuf {
		d.Buffer.Destroy()
		d.Buffer = nil
		d.ownedbuf = false
	}
}

func (d *Uniform[K]) Set(data K) {
	d.Buffer.Set(data)
}

func (d *Uniform[K]) write() {
	d.set.Write(core1_0.WriteDescriptorSet{
		DstBinding:      d.binding,
		DstArrayElement: 0,
		DescriptorType:  core1_0.DescriptorTypeUniformBuffer,
		BufferInfo: []core1_0.DescriptorBufferInfo{
			{
				Buffer: d.Buffer.Ptr(),
				Offset: 0,
				Range:  d.Buffer.Size(),
			},
		},
	})
}

func (d *Uniform[K]) LayoutBinding(binding int) core1_0.DescriptorSetLayoutBinding {
	return core1_0.DescriptorSetLayoutBinding{
		Binding:         binding,
		DescriptorType:  core1_0.DescriptorTypeUniformBuffer,
		DescriptorCount: 1,
		StageFlags:      core1_0.ShaderStageFlags(d.Stages),
	}
}

func (d *Uniform[K]) BindingFlags() ext_descriptor_indexing.DescriptorBindingFlags { return 0 }
