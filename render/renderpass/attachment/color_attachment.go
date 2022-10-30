package attachment

import (
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/device"

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
	Usage         vk.ImageUsageFlagBits

	// Allocation strategy. Defaults to allocating new images.
	Allocator Allocator

	Blend Blend
}

func (desc *Color) defaults() {
	if desc.Samples == 0 {
		desc.Samples = vk.SampleCount1Bit
	}
	if desc.Allocator == nil {
		desc.Allocator = &alloc{}
	}
}

func NewColor(device device.T, desc Color) T {
	desc.defaults()

	var clear vk.ClearValue
	clear.SetColor([]float32{desc.Clear.R, desc.Clear.G, desc.Clear.B, desc.Clear.A})

	return &attachment{
		name:   desc.Name,
		alloc:  desc.Allocator,
		clear:  clear,
		blend:  desc.Blend,
		format: desc.Format,
		usage:  desc.Usage,
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
