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

func NewLayout(device device.T, descriptors []descriptor.T) Layout {
	info := vk.PipelineLayoutCreateInfo{
		SType:          vk.StructureTypePipelineLayoutCreateInfo,
		SetLayoutCount: uint32(len(descriptors)),
		PSetLayouts: util.Map(descriptors, func(i int, desc descriptor.T) vk.DescriptorSetLayout {
			return desc.Ptr()
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
	vk.DestroyPipelineLayout(l.device.Ptr(), l.ptr, nil)
	l.ptr = nil
}
