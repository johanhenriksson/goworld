package descriptor

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/ext_descriptor_indexing"
)

type Sampler struct {
	Stages  core1_0.ShaderStageFlags
	Texture *texture.Texture

	binding int
	sampler core1_0.Sampler
	view    core1_0.ImageView
	set     Set
}

var _ Descriptor = (*Sampler)(nil)

func (d *Sampler) Initialize(device *device.Device, set Set, binding int) {
	d.set = set
	d.binding = binding

	if d.Texture != nil {
		d.Set(d.Texture)
	}
}

func (d *Sampler) String() string {
	return fmt.Sprintf("Sampler:%d", d.binding)
}

func (d *Sampler) Destroy() {}

func (d *Sampler) Set(tex *texture.Texture) {
	d.Texture = tex
	d.sampler = tex.Ptr()
	d.view = tex.View().Ptr()
	d.write()
}

func (d *Sampler) LayoutBinding(binding int) core1_0.DescriptorSetLayoutBinding {
	return core1_0.DescriptorSetLayoutBinding{
		Binding:         binding,
		DescriptorType:  core1_0.DescriptorTypeCombinedImageSampler,
		DescriptorCount: 1,
		StageFlags:      core1_0.ShaderStageFlags(d.Stages),
	}
}

func (d *Sampler) BindingFlags() ext_descriptor_indexing.DescriptorBindingFlags { return 0 }

func (d *Sampler) write() {
	d.set.Write(core1_0.WriteDescriptorSet{
		DstSet:          d.set.Ptr(),
		DstBinding:      d.binding,
		DstArrayElement: 0,
		DescriptorType:  core1_0.DescriptorTypeCombinedImageSampler,
		ImageInfo: []core1_0.DescriptorImageInfo{
			{
				Sampler:     d.sampler,
				ImageView:   d.view,
				ImageLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
			},
		},
	})
}
