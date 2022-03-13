package renderpass

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	vk "github.com/vulkan-go/vulkan"
)

type Attachment interface {
	Name() string
	Image(int) image.T
	View(int) image.View
	Clear() vk.ClearValue
	Description() vk.AttachmentDescription
	Destroy()
	Blend() bool
}

type attachment struct {
	name     string
	view     []image.View
	image    []image.T
	clear    vk.ClearValue
	desc     vk.AttachmentDescription
	imgowner bool
	blend    bool
}

func (a *attachment) Description() vk.AttachmentDescription {
	return a.desc
}

func (a *attachment) Name() string              { return a.name }
func (a *attachment) Image(frame int) image.T   { return a.image[frame] }
func (a *attachment) View(frame int) image.View { return a.view[frame] }
func (a *attachment) Clear() vk.ClearValue      { return a.clear }
func (a *attachment) Blend() bool               { return a.blend }

func (a *attachment) Destroy() {
	for i := range a.image {
		a.view[i].Destroy()
		if a.imgowner {
			a.image[i].Destroy()
		}
	}
}

func NewColorAttachment(device device.T, desc ColorAttachment, frames, width, height int) Attachment {
	desc.defaults()

	images := desc.Images
	imgowner := false
	if len(images) == 0 {
		imgowner = true
		images = make([]image.T, frames)
		for i := range images {
			images[i] = image.New2D(
				device,
				width, height, desc.Format,
				vk.ImageUsageFlags(vk.ImageUsageColorAttachmentBit|desc.Usage))
		}
	} else if len(images) != frames {
		panic("wrong number of images supplied")
	}

	views := make([]image.View, frames)
	for i := range views {
		views[i] = images[i].View(desc.Format, vk.ImageAspectFlags(vk.ImageAspectColorBit))
	}

	var clear vk.ClearValue
	clear.SetColor([]float32{desc.Clear.R, desc.Clear.G, desc.Clear.B, desc.Clear.A})

	return &attachment{
		name:     desc.Name,
		image:    images,
		view:     views,
		clear:    clear,
		imgowner: imgowner,
		blend:    desc.Blend,
		desc: vk.AttachmentDescription{
			Format:        desc.Format,
			Samples:       desc.Samples,
			LoadOp:        desc.LoadOp,
			StoreOp:       desc.StoreOp,
			InitialLayout: desc.InitialLayout,
			FinalLayout:   desc.FinalLayout,

			// color attachments dont have stencil buffers, so we dont care about them
			StencilLoadOp:  vk.AttachmentLoadOpDontCare,
			StencilStoreOp: vk.AttachmentStoreOpDontCare,
		},
	}
}

func NewDepthAttachment(device device.T, desc DepthAttachment, frames, width, height, index int) Attachment {
	desc.defaults()

	depthFormat := device.GetDepthFormat()

	images := desc.Images
	imgowner := false
	if len(images) == 0 {
		imgowner = true
		images = make([]image.T, frames)
		for i := range images {
			images[i] = image.New2D(
				device, width, height, depthFormat,
				vk.ImageUsageFlags(vk.ImageUsageDepthStencilAttachmentBit|desc.Usage))
		}
	} else if len(images) != frames {
		panic("wrong number of images supplied")
	}

	views := make([]image.View, frames)
	for i := range images {
		views[i] = images[i].View(depthFormat, vk.ImageAspectFlags(vk.ImageAspectDepthBit))
	}

	var clear vk.ClearValue
	clear.SetDepthStencil(desc.ClearDepth, desc.ClearStencil)

	return &attachment{
		name:     "depth",
		image:    images,
		view:     views,
		clear:    clear,
		imgowner: imgowner,
		blend:    false,
		desc: vk.AttachmentDescription{
			Format:         depthFormat,
			Samples:        desc.Samples,
			LoadOp:         desc.LoadOp,
			StoreOp:        desc.StoreOp,
			StencilLoadOp:  desc.StencilLoadOp,
			StencilStoreOp: desc.StencilStoreOp,
			InitialLayout:  desc.InitialLayout,
			FinalLayout:    desc.FinalLayout,
		},
	}
}
