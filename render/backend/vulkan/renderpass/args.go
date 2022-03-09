package renderpass

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/color"
	vk "github.com/vulkan-go/vulkan"
)

type Args struct {
	Frames int
	Width  int
	Height int

	ColorAttachments []ColorAttachment
	DepthAttachment  *DepthAttachment

	Subpasses    []Subpass
	Dependencies []SubpassDependency
}

type ColorAttachment struct {
	Name          string
	Format        vk.Format
	Samples       vk.SampleCountFlagBits
	LoadOp        vk.AttachmentLoadOp
	StoreOp       vk.AttachmentStoreOp
	InitialLayout vk.ImageLayout
	FinalLayout   vk.ImageLayout
	Clear         color.T
	Images        []image.T
}

func (desc *ColorAttachment) defaults() {
	if desc.Samples == 0 {
		desc.Samples = vk.SampleCount1Bit
	}
}

type DepthAttachment struct {
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
}

func (desc *DepthAttachment) defaults() {
	if desc.Samples == 0 {
		desc.Samples = vk.SampleCount1Bit
	}
}

type Subpass struct {
	Name             string
	Depth            bool
	ColorAttachments []string
}

type SubpassDependency struct {
	Src string
	Dst string

	Flags         vk.DependencyFlags
	SrcStageMask  vk.PipelineStageFlags
	SrcAccessMask vk.AccessFlags
	DstStageMask  vk.PipelineStageFlags
	DstAccessMask vk.AccessFlags
}
