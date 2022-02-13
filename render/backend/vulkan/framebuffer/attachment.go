package framebuffer

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"

	vk "github.com/vulkan-go/vulkan"
)

type Attachment interface {
	IsDepth() bool
}

type attach struct {
	name      string
	view      image.View
	desc      vk.AttachmentDescription
	ref       vk.AttachmentReference
	inlayout  vk.ImageLayout
	outlayout vk.ImageLayout
}

type RenderPassArgs struct {
	Flags        vk.RenderPassCreateFlags
	Attachments  []vk.AttachmentDescription
	Subpasses    []vk.SubpassDescription
	Dependencies []vk.SubpassDependency
}

func NewAt(device device.T, args RenderPassArgs) {
	info := vk.RenderPassCreateInfo{
		SType:           vk.StructureTypeRenderPassCreateInfo,
		Flags:           args.Flags,
		AttachmentCount: uint32(len(args.Attachments)),
		PAttachments:    args.Attachments,
		SubpassCount:    uint32(len(args.Subpasses)),
		PSubpasses:      args.Subpasses,
		DependencyCount: uint32(len(args.Dependencies)),
		PDependencies:   args.Dependencies,
	}

	var ptr vk.RenderPass
	vk.CreateRenderPass(device.Ptr(), &info, nil, &ptr)
}
