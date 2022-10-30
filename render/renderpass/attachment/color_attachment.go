package attachment

import (
	"log"

	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"

	vk "github.com/vulkan-go/vulkan"
)

type Color struct {
	Name          Name
	Format        vk.Format
	Samples       vk.SampleCountFlagBits
	LoadOp        vk.AttachmentLoadOp
	StoreOp       vk.AttachmentStoreOp
	InitialLayout vk.ImageLayout
	FinalLayout   vk.ImageLayout
	Clear         color.T
	Images        []image.T
	Usage         vk.ImageUsageFlagBits

	Blend Blend
}

func (desc *Color) defaults() {
	if desc.Samples == 0 {
		desc.Samples = vk.SampleCount1Bit
	}
}

func NewColor(device device.T, desc Color, frames, width, height int) T {
	desc.defaults()

	images := desc.Images
	imgowner := false
	if len(images) == 0 {
		log.Println("  allocating", frames, "color attachments")
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
