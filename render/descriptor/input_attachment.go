package descriptor

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type InputAttachment struct {
	Stages core1_0.ShaderStageFlags
	Layout core1_0.ImageLayout

	binding int
	view    core1_0.ImageView
	set     Set
}

var _ Descriptor = &InputAttachment{}

func (d *InputAttachment) Initialize(device *device.Device) {
	if d.Layout == 0 {
		d.Layout = core1_0.ImageLayoutShaderReadOnlyOptimal
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

func (d *InputAttachment) LayoutBinding(binding int) core1_0.DescriptorSetLayoutBinding {
	d.binding = binding
	return core1_0.DescriptorSetLayoutBinding{
		Binding:         binding,
		DescriptorType:  core1_0.DescriptorTypeInputAttachment,
		DescriptorCount: 1,
		StageFlags:      core1_0.ShaderStageFlags(d.Stages),
	}
}

func (d *InputAttachment) BindingFlags() ext_descriptor_indexing.DescriptorBindingFlags { return 0 }

func (d *InputAttachment) write() {
	d.set.Write(core1_0.WriteDescriptorSet{
		DstSet:          d.set.Ptr(),
		DstBinding:      d.binding,
		DstArrayElement: 0,
		DescriptorType:  core1_0.DescriptorTypeInputAttachment,
		ImageInfo: []core1_0.DescriptorImageInfo{
			{
				ImageView:   d.view,
				ImageLayout: d.Layout,
			},
		},
	})
}
