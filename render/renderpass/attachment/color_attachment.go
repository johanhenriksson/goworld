package attachment

import (
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Color struct {
	Name          Name
	Format        core1_0.Format
	Samples       core1_0.SampleCountFlags
	LoadOp        core1_0.AttachmentLoadOp
	StoreOp       core1_0.AttachmentStoreOp
	InitialLayout core1_0.ImageLayout
	FinalLayout   core1_0.ImageLayout
	Clear         color.T
	Usage         core1_0.ImageUsageFlags

	// Allocation strategy. Defaults to allocating new images.
	Allocator Allocator

	Blend Blend
}

func (desc *Color) defaults() {
	if desc.Samples == 0 {
		desc.Samples = core1_0.Samples1
	}
	if desc.Allocator == nil {
		desc.Allocator = &alloc{}
	}
}

func NewColor(device device.T, desc Color) T {
	desc.defaults()

	clear := core1_0.ClearValueFloat{desc.Clear.R, desc.Clear.G, desc.Clear.B, desc.Clear.A}

	return &attachment{
		name:   desc.Name,
		alloc:  desc.Allocator,
		clear:  clear,
		blend:  desc.Blend,
		format: desc.Format,
		usage:  desc.Usage,
		desc: core1_0.AttachmentDescription{
			Format:        desc.Format,
			Samples:       desc.Samples,
			LoadOp:        desc.LoadOp,
			StoreOp:       desc.StoreOp,
			InitialLayout: desc.InitialLayout,
			FinalLayout:   desc.FinalLayout,

			// color attachments dont have stencil buffers, so we dont care about them
			StencilLoadOp:  core1_0.AttachmentLoadOpDontCare,
			StencilStoreOp: core1_0.AttachmentStoreOpDontCare,
		},
	}
}
