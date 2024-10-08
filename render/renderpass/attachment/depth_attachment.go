package attachment

import (
	"github.com/johanhenriksson/goworld/render/device"

	"github.com/vkngwrapper/core/v2/core1_0"
)

const DepthName Name = "depth"

type Depth struct {
	Samples        core1_0.SampleCountFlags
	LoadOp         core1_0.AttachmentLoadOp
	StoreOp        core1_0.AttachmentStoreOp
	StencilLoadOp  core1_0.AttachmentLoadOp
	StencilStoreOp core1_0.AttachmentStoreOp
	InitialLayout  core1_0.ImageLayout
	FinalLayout    core1_0.ImageLayout
	ClearDepth     float32
	ClearStencil   uint32

	// Allocation strategy. Defaults to allocating new images.
	Image Image
}

func (desc *Depth) defaults() {
	if desc.Samples == 0 {
		desc.Samples = core1_0.Samples1
	}
	if desc.Image == nil {
		panic("no image reference")
	}
}

func NewDepth(device *device.Device, desc Depth) T {
	desc.defaults()

	clear := core1_0.ClearValueDepthStencil{
		Depth:   desc.ClearDepth,
		Stencil: desc.ClearStencil,
	}

	return &attachment{
		name:  DepthName,
		image: desc.Image,
		clear: clear,
		desc: core1_0.AttachmentDescription{
			Format:         desc.Image.Format(),
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
