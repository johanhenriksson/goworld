package pipeline

import (
	"log"

	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/util"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Layout struct {
	ptr    core1_0.PipelineLayout
	device *device.Device
}

func NewLayout(device *device.Device, descriptors []descriptor.SetLayout, constants []PushConstant) *Layout {
	offset := 0
	info := core1_0.PipelineLayoutCreateInfo{

		SetLayouts: util.Map(descriptors, func(desc descriptor.SetLayout) core1_0.DescriptorSetLayout {
			return desc.Ptr()
		}),

		PushConstantRanges: util.Map(constants, func(push PushConstant) core1_0.PushConstantRange {
			size := push.Size()
			log.Printf("push: %d bytes", size)
			rng := core1_0.PushConstantRange{
				StageFlags: core1_0.ShaderStageFlags(push.Stages),
				Offset:     offset,
				Size:       size,
			}
			offset += size
			return rng
		}),
	}

	ptr, _, err := device.Ptr().CreatePipelineLayout(nil, info)
	if err != nil {
		panic(err)
	}

	return &Layout{
		ptr:    ptr,
		device: device,
	}
}

func (l *Layout) Ptr() core1_0.PipelineLayout {
	return l.ptr
}

func (l *Layout) Destroy() {
	if l.ptr != nil {
		l.ptr.Destroy(nil)
		l.ptr = nil
	}
}
