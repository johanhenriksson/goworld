package pipeline

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/descriptor"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/util"
	vk "github.com/vulkan-go/vulkan"
)

type Layout interface {
	device.Resource[vk.PipelineLayout]
}

type layout struct {
	ptr    vk.PipelineLayout
	device device.T
}

func NewLayout(device device.T, descriptors []descriptor.SetLayout, constants []PushConstant) Layout {
	info := vk.PipelineLayoutCreateInfo{
		SType: vk.StructureTypePipelineLayoutCreateInfo,

		SetLayoutCount: uint32(len(descriptors)),
		PSetLayouts: util.Map(descriptors, func(desc descriptor.SetLayout) vk.DescriptorSetLayout {
			return desc.Ptr()
		}),

		PushConstantRangeCount: uint32(len(constants)),
		PPushConstantRanges: util.Map(constants, func(push PushConstant) vk.PushConstantRange {
			return vk.PushConstantRange{
				StageFlags: vk.ShaderStageFlags(push.Stages),
				Offset:     uint32(push.Offset),
				Size:       uint32(push.Size),
			}
		}),
	}

	var ptr vk.PipelineLayout
	vk.CreatePipelineLayout(device.Ptr(), &info, nil, &ptr)

	return &layout{
		ptr:    ptr,
		device: device,
	}
}

func (l *layout) Ptr() vk.PipelineLayout {
	return l.ptr
}

func (l *layout) Destroy() {
	if l.ptr != nil {
		vk.DestroyPipelineLayout(l.device.Ptr(), l.ptr, nil)
		l.ptr = nil
	}
}
