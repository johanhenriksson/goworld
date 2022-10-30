package attachment

import (
	"log"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	vk "github.com/vulkan-go/vulkan"
)

const DepthName Name = "depth"

type Depth struct {
	Samples        vk.SampleCountFlagBits
	LoadOp         vk.AttachmentLoadOp
	StoreOp        vk.AttachmentStoreOp
	StencilLoadOp  vk.AttachmentLoadOp
	StencilStoreOp vk.AttachmentStoreOp
	InitialLayout  vk.ImageLayout
	FinalLayout    vk.ImageLayout
	ClearDepth     float32
	ClearStencil   uint32
	Images         []image.T
	Usage          vk.ImageUsageFlagBits
}

func (desc *Depth) defaults() {
	if desc.Samples == 0 {
		desc.Samples = vk.SampleCount1Bit
	}
}

func NewDepth(device device.T, desc Depth, frames, width, height, index int) T {
	desc.defaults()

	depthFormat := device.GetDepthFormat()

	images := desc.Images
	imgowner := false
	if len(images) == 0 {
		log.Println("  allocating", frames, "depth attachments")
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
		name:     DepthName,
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
