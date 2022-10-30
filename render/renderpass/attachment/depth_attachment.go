package attachment

import (
	"github.com/johanhenriksson/goworld/render/device"
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
	Usage          vk.ImageUsageFlagBits
	ClearDepth     float32
	ClearStencil   uint32

	// Allocation strategy. Defaults to allocating new images.
	Allocator Allocator
}

func (desc *Depth) defaults() {
	if desc.Samples == 0 {
		desc.Samples = vk.SampleCount1Bit
	}
	if desc.Allocator == nil {
		desc.Allocator = &alloc{}
	}
}

func NewDepth(device device.T, desc Depth) T {
	desc.defaults()

	depthFormat := device.GetDepthFormat()

	var clear vk.ClearValue
	clear.SetDepthStencil(desc.ClearDepth, desc.ClearStencil)

	return &attachment{
		name:   DepthName,
		alloc:  desc.Allocator,
		clear:  clear,
		format: depthFormat,
		usage:  desc.Usage,
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
