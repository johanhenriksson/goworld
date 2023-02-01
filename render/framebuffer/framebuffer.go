package framebuffer

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vkerror"
	"github.com/johanhenriksson/goworld/util"

	"github.com/vkngwrapper/core/v2/core1_0"
	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[core1_0.Framebuffer]

	Attachment(attachment.Name) image.View
	Size() (int, int)
}

type framebuf struct {
	ptr         core1_0.Framebuffer
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

	allocate := func(attach attachment.T, usage core1_0.ImageUsageFlags, aspect core1_0.ImageAspectFlags) error {
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

		view, err := img.View(attach.Format(), core1_0.ImageAspectFlags(aspect))
		if err != nil {
			return err
		}
		views = append(views, view)

		attachs[attach.Name()] = view
		return nil
	}

	for _, attach := range attachments {
		if err := allocate(attach, core1_0.ImageUsageColorAttachment, core1_0.ImageAspectColor); err != nil {
			cleanup()
			return nil, err
		}
	}
	if depth != nil {
		if err := allocate(depth, core1_0.ImageUsageDepthStencilAttachment, core1_0.ImageAspectDepth); err != nil {
			cleanup()
			return nil, err
		}
	}

	info := core1_0.FramebufferCreateInfo{
		RenderPass:  pass.Ptr(),
		Attachments: util.Map(views, func(v image.View) core1_0.ImageView { return v.Ptr() }),
		Width:       width,
		Height:      height,
		Layers:      1,
	}

	var ptr core1_0.Framebuffer
	ptr, result, err := device.Ptr().CreateFramebuffer(nil, info)
	if err != nil {
		panic(err)
	}
	if result != core1_0.VKSuccess {
		cleanup()
		return nil, vkerror.FromResult(result)
	}

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

func (b *framebuf) Ptr() core1_0.Framebuffer {
	return b.ptr
}

func (b *framebuf) Size() (int, int) {
	return b.width, b.height
}

func (b *framebuf) Attachment(name attachment.Name) image.View {
	return b.attachments[name]
}

func (b *framebuf) Destroy() {
	if b.ptr == nil {
		panic("framebuffer already destroyed")
	}
	for _, view := range b.views {
		view.Destroy()
	}
	b.views = nil
	for _, image := range b.images {
		image.Destroy()
	}
	b.images = nil
	b.attachments = nil

	b.ptr.Destroy(nil)
	b.ptr = nil
	b.device = nil
}
