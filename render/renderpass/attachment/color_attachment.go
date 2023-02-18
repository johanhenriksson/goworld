package attachment

import (
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Color struct {
	Name          Name
	Samples       core1_0.SampleCountFlags
	LoadOp        core1_0.AttachmentLoadOp
	StoreOp       core1_0.AttachmentStoreOp
	InitialLayout core1_0.ImageLayout
	FinalLayout   core1_0.ImageLayout
	Clear         color.T
	Image         Image
	Blend         Blend
}

func (desc *Color) defaults() {
	if desc.Samples == 0 {
		desc.Samples = core1_0.Samples1
	}
	if desc.Image == nil {
		panic("no image reference")
	}
}

func NewColor(device device.T, desc Color) T {
	desc.defaults()

	clear := core1_0.ClearValueFloat{desc.Clear.R, desc.Clear.G, desc.Clear.B, desc.Clear.A}

	return &attachment{
		name:  desc.Name,
		image: desc.Image,
		clear: clear,
		blend: desc.Blend,
		desc: core1_0.AttachmentDescription{
			Format:        desc.Image.Format(),
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
