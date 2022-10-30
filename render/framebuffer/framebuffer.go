package framebuffer

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Framebuffer]

	Attachment(attachment.Name) image.View
	Size() (int, int)
}

type framebuf struct {
	ptr         vk.Framebuffer
	device      device.T
	attachments map[attachment.Name]image.View
	views       []image.View
	images      []image.T
	width       int
	height      int
}

func New(device device.T, width, height int, pass renderpass.T) (T, error) {
	attachments := pass.Attachments()
	depth := pass.Depth()

	images := make([]image.T, 0, len(attachments)+1)
	views := make([]image.View, 0, len(attachments)+1)
	attachs := make(map[attachment.Name]image.View)

	cleanup := func() {
		// clean up the mess we've made so far
		for _, view := range views {
			view.Destroy()
		}
		for _, image := range images {
			image.Destroy()
		}
	}

	allocate := func(attach attachment.T, usage vk.ImageUsageFlagBits, aspect vk.ImageAspectFlagBits) error {
		img, err := attach.Allocator().Alloc(
			device,
			width, height,
			attach.Format(),
			usage|attach.Usage(),
		)
		if err != nil {
			return err
		}
		images = append(images, img)

		view := img.View(attach.Format(), vk.ImageAspectFlags(aspect))
		views = append(views, view)

		attachs[attach.Name()] = view
		return nil
	}

	for _, attach := range attachments {
		if err := allocate(attach, vk.ImageUsageColorAttachmentBit, vk.ImageAspectColorBit); err != nil {
			cleanup()
			return nil, err
		}
	}
	if depth != nil {
		if err := allocate(depth, vk.ImageUsageDepthStencilAttachmentBit, vk.ImageAspectDepthBit); err != nil {
			cleanup()
			return nil, err
		}
	}

	info := vk.FramebufferCreateInfo{
		SType:           vk.StructureTypeFramebufferCreateInfo,
		RenderPass:      pass.Ptr(),
		AttachmentCount: uint32(len(attachs)),
		PAttachments:    util.Map(views, func(v image.View) vk.ImageView { return v.Ptr() }),
		Width:           uint32(width),
		Height:          uint32(height),
		Layers:          1,
	}

	var ptr vk.Framebuffer
	vk.CreateFramebuffer(device.Ptr(), &info, nil, &ptr)

	return &framebuf{
		ptr:         ptr,
		device:      device,
		width:       width,
		height:      height,
		images:      images,
		views:       views,
		attachments: attachs,
	}, nil
}

func (b *framebuf) Ptr() vk.Framebuffer {
	return b.ptr
}

func (b *framebuf) Size() (int, int) {
	return b.width, b.height
}

func (b *framebuf) Attachment(name attachment.Name) image.View {
	return b.attachments[name]
}

func (b *framebuf) Destroy() {
	for _, view := range b.views {
		view.Destroy()
	}
	b.views = nil
	for _, image := range b.images {
		image.Destroy()
	}
	b.images = nil
	b.attachments = nil

	vk.DestroyFramebuffer(b.device.Ptr(), b.ptr, nil)
	b.ptr = nil
	b.device = nil
}
