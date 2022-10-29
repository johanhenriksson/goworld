package descriptor

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"

	vk "github.com/vulkan-go/vulkan"
)

type InputAttachment struct {
	Stages vk.ShaderStageFlagBits
	Layout vk.ImageLayout

	binding int
	view    vk.ImageView
	set     Set
}

var _ Descriptor = &InputAttachment{}

func (d *InputAttachment) Initialize(device device.T) {
	if d.Layout == 0 {
		d.Layout = vk.ImageLayoutShaderReadOnlyOptimal
	}
}

func (d *InputAttachment) String() string {
	return fmt.Sprintf("Input:%d", d.binding)
}

func (d *InputAttachment) Destroy() {}

func (d *InputAttachment) Bind(set Set, binding int) {
	d.set = set
	d.binding = binding
}

func (d *InputAttachment) Set(view image.View) {
	d.view = view.Ptr()
	d.write()
}

func (d *InputAttachment) LayoutBinding(binding int) vk.DescriptorSetLayoutBinding {
	d.binding = binding
	return vk.DescriptorSetLayoutBinding{
		Binding:         uint32(binding),
		DescriptorType:  vk.DescriptorTypeInputAttachment,
		DescriptorCount: 1,
		StageFlags:      vk.ShaderStageFlags(d.Stages),
	}
}

func (d *InputAttachment) BindingFlags() vk.DescriptorBindingFlags { return 0 }

func (d *InputAttachment) write() {
	d.set.Write(vk.WriteDescriptorSet{
		SType:           vk.StructureTypeWriteDescriptorSet,
		DstSet:          d.set.Ptr(),
		DstBinding:      uint32(d.binding),
		DstArrayElement: 0,
		DescriptorCount: 1,
		DescriptorType:  vk.DescriptorTypeInputAttachment,
		PImageInfo: []vk.DescriptorImageInfo{
			{
				ImageView:   d.view,
				ImageLayout: d.Layout,
			},
		},
	})
}
