package renderpass

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	vk "github.com/vulkan-go/vulkan"
)

type Attachment interface {
	Index() int
	Image(int) image.T
	View(int) image.View
	Clear() vk.ClearValue
	Description() vk.AttachmentDescription
	Destroy()
}

type attachment struct {
	index    int
	view     []image.View
	image    []image.T
	clear    vk.ClearValue
	desc     vk.AttachmentDescription
	imgowner bool
}

func (a *attachment) Description() vk.AttachmentDescription {
	return a.desc
}

func (a *attachment) Index() int                { return a.index }
func (a *attachment) Image(frame int) image.T   { return a.image[frame] }
func (a *attachment) View(frame int) image.View { return a.view[frame] }
func (a *attachment) Clear() vk.ClearValue      { return a.clear }

func (a *attachment) Destroy() {
	for i := range a.image {
		a.view[i].Destroy()
		if a.imgowner {
			a.image[i].Destroy()
		}
	}
}

func NewColorAttachment(device device.T, desc ColorAttachment, frames, width, height int) Attachment {
	images := desc.Images
	imgowner := false
	if len(images) == 0 {
		imgowner = true
		images = make([]image.T, frames)
		for i := range images {
			images[i] = image.New2D(
				device,
				width, height, desc.Format,
				vk.ImageUsageFlags(vk.ImageUsageColorAttachmentBit|vk.ImageUsageSampledBit))
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
		index:    desc.Index,
		image:    images,
		view:     views,
		clear:    clear,
		imgowner: imgowner,
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
	depthFormat := device.GetDepthFormat()
	depthUsage := vk.ImageUsageFlags(vk.ImageUsageDepthStencilAttachmentBit | vk.ImageUsageSampledBit)

	images := desc.Images
	imgowner := false
	if len(images) == 0 {
		imgowner = true
		images = make([]image.T, frames)
		for i := range images {
			images[i] = image.New2D(device, width, height, depthFormat, depthUsage)
		}
	} else if len(images) != frames {
		panic("wrong number of images supplied")
	}

	views := make([]image.View, frames)
	for i := range images {
		views[i] = images[i].View(depthFormat, vk.ImageAspectFlags(vk.ImageAspectDepthBit|vk.ImageAspectStencilBit))
	}

	var clear vk.ClearValue
	clear.SetDepthStencil(desc.ClearDepth, desc.ClearStencil)

	return &attachment{
		index:    index,
		image:    images,
		view:     views,
		clear:    clear,
		imgowner: imgowner,
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
