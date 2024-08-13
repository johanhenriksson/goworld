package pipeline

import (
	"log"

	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/samber/lo"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type Layout struct {
	ptr    core1_0.PipelineLayout
	device *device.Device
}

func NewLayout(device *device.Device, descriptors []descriptor.SetLayout, constants []PushConstant) *Layout {
	offset := 0
	info := core1_0.PipelineLayoutCreateInfo{

		SetLayouts: lo.Map(descriptors, func(desc descriptor.SetLayout, _ int) core1_0.DescriptorSetLayout {
			return desc.Ptr()
		}),

		PushConstantRanges: lo.Map(constants, func(push PushConstant, _ int) core1_0.PushConstantRange {
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
