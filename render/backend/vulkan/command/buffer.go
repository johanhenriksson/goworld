package command

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"

	vk "github.com/vulkan-go/vulkan"
)

type Buffer interface {
	Ptr() vk.CommandBuffer
}

type buffer struct {
	ptr    vk.CommandBuffer
	device device.T
}

func newBuffer(device device.T, ptr vk.CommandBuffer) Buffer {
	return &buffer{
		ptr:    ptr,
		device: device,
	}
}

func (b *buffer) Ptr() vk.CommandBuffer {
	return b.ptr
}
