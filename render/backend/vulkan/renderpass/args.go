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
	Usage         vk.ImageUsageFlagBits

	Blend Blend
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
	Usage          vk.ImageUsageFlagBits
}

func (desc *DepthAttachment) defaults() {
	if desc.Samples == 0 {
		desc.Samples = vk.SampleCount1Bit
	}
}

type Subpass struct {
	index            int
	Name             string
	Depth            bool
	ColorAttachments []string
	InputAttachments []string
}

func (s *Subpass) Index() int {
	return s.index
}

type SubpassDependency struct {
	Src string
	Dst string

	Flags         vk.DependencyFlagBits
	SrcStageMask  vk.PipelineStageFlagBits
	SrcAccessMask vk.AccessFlagBits
	DstStageMask  vk.PipelineStageFlagBits
	DstAccessMask vk.AccessFlagBits
}
