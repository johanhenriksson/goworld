package framebuffer

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/renderpass"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/render/vkerror"

	"github.com/samber/lo"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type Framebuffer struct {
	ptr         core1_0.Framebuffer
	device      *device.Device
	name        string
	attachments map[attachment.Name]*image.View
	views       image.ViewArray
	images      image.Array
	width       int
	height      int
}

func New(device *device.Device, name string, width, height int, pass *renderpass.Renderpass) (*Framebuffer, error) {
	attachments := pass.Attachments()
	depth := pass.Depth()

	images := make(image.Array, 0, len(attachments)+1)
	views := make(image.ViewArray, 0, len(attachments)+1)
	attachs := make(map[attachment.Name]*image.View)

	cleanup := func() {
		// clean up the mess we've made so far
		for _, view := range views {
			view.Destroy()
		}
		for _, image := range images {
			image.Destroy()
		}
	}

	allocate := func(attach attachment.T, aspect core1_0.ImageAspectFlags) error {
		img, ownership, err := attach.Image().Next(
			device,
			name,
			width, height,
		)
		if err != nil {
			return err
		}
		if ownership {
			// the framebuffer is responsible for deallocating the image
			images = append(images, img)
		}

		view, err := img.View(img.Format(), core1_0.ImageAspectFlags(aspect))
		if err != nil {
			return err
		}
		views = append(views, view)

		attachs[attach.Name()] = view
		return nil
	}

	for _, attach := range attachments {
		if err := allocate(attach, core1_0.ImageAspectColor); err != nil {
			cleanup()
			return nil, err
		}
	}
	if depth != nil {
		if err := allocate(depth, core1_0.ImageAspectDepth); err != nil {
			cleanup()
			return nil, err
		}
	}

	info := core1_0.FramebufferCreateInfo{
		RenderPass:  pass.Ptr(),
		Attachments: lo.Map(views, func(v *image.View, _ int) core1_0.ImageView { return v.Ptr() }),
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

	device.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()), core1_0.ObjectTypeFramebuffer, name)

	return &Framebuffer{
		ptr:         ptr,
		device:      device,
		name:        name,
		width:       width,
		height:      height,
		images:      images,
		views:       views,
		attachments: attachs,
	}, nil
}

func (b *Framebuffer) Ptr() core1_0.Framebuffer {
	return b.ptr
}

func (b *Framebuffer) Size() (int, int) {
	return b.width, b.height
}

func (b *Framebuffer) Attachment(name attachment.Name) *image.View {
	return b.attachments[name]
}

func (b *Framebuffer) Destroy() {
	if b.ptr == nil {
		panic("framebuffer already destroyed")
	}

	for _, image := range b.images {
		image.Destroy()
	}
	b.images = nil

	for _, view := range b.views {
		view.Destroy()
	}
	b.views = nil

	b.attachments = nil

	b.ptr.Destroy(nil)
	b.ptr = nil
	b.device = nil
}
