package framebuffer

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Framebuffer]

	Attachments() []image.View
	Size() (int, int)
}

type framebuf struct {
	ptr         vk.Framebuffer
	device      device.T
	attachments []image.View
	width       int
	height      int
}

func New(device device.T, width, height int, pass vk.RenderPass, attachments []image.View) T {
	info := vk.FramebufferCreateInfo{
		SType:           vk.StructureTypeFramebufferCreateInfo,
		RenderPass:      pass,
		AttachmentCount: uint32(len(attachments)),
		PAttachments:    util.Map(attachments, func(v image.View) vk.ImageView { return v.Ptr() }),
		Width:           uint32(width),
		Height:          uint32(height),
		Layers:          1,
	}

	var ptr vk.Framebuffer
	vk.CreateFramebuffer(device.Ptr(), &info, nil, &ptr)

	return &framebuf{
		ptr:         ptr,
		device:      device,
		attachments: attachments,
		width:       width,
		height:      height,
	}
}

func (b *framebuf) Ptr() vk.Framebuffer {
	return b.ptr
}

func (b *framebuf) Size() (int, int) {
	return b.width, b.height
}

func (b *framebuf) Attachments() []image.View {
	return b.attachments
}

func (b *framebuf) Destroy() {
	vk.DestroyFramebuffer(b.device.Ptr(), b.ptr, nil)
	b.ptr = nil
}
