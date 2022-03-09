package renderpass

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/framebuffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.RenderPass]

	Depth() Attachment
	Attachment(name string) Attachment
	Attachments() []Attachment
	Framebuffer(frame int) framebuffer.T
	Clear() []vk.ClearValue
}

type renderpass struct {
	device       device.T
	ptr          vk.RenderPass
	framebuffers []framebuffer.T
	attachments  map[string]Attachment
	depth        Attachment
	order        []string
	clear        []vk.ClearValue
}

func New(device device.T, args Args) T {
	clear := make([]vk.ClearValue, 0, len(args.ColorAttachments)+1)
	attachments := make(map[string]Attachment)
	for name, desc := range args.ColorAttachments {
		attachment := NewColorAttachment(device, desc, args.Frames, args.Width, args.Height)
		clear = append(clear, attachment.Clear())
		attachments[name] = attachment
	}

	var depth Attachment
	if args.DepthAttachment != nil {
		index := len(attachments)
		depth = NewDepthAttachment(device, *args.DepthAttachment, args.Frames, args.Width, args.Height, index)
		clear = append(clear, depth.Clear())
	}

	descriptions := make([]vk.AttachmentDescription, 0, len(args.ColorAttachments)+1)
	attachmentOrder := make([]string, 0, len(args.ColorAttachments))
	for name, attachment := range attachments {
		descriptions = append(descriptions, attachment.Description())
		attachmentOrder = append(attachmentOrder, name)
	}
	if depth != nil {
		descriptions = append(descriptions, depth.Description())
	}

	subpasses := make([]vk.SubpassDescription, 0, len(args.Subpasses))
	subpassIndices := make(map[string]int)
	for idx, subpass := range args.Subpasses {
		var depthRef *vk.AttachmentReference
		if depth != nil {
			depthRef = &vk.AttachmentReference{
				Attachment: uint32(len(attachmentOrder)),
				Layout:     vk.ImageLayoutDepthStencilAttachmentOptimal,
			}
		}
		subpasses = append(subpasses, vk.SubpassDescription{
			PipelineBindPoint:    vk.PipelineBindPointGraphics,
			ColorAttachmentCount: uint32(len(attachments)),
			PColorAttachments: util.Map(attachmentOrder, func(idx int, name string) vk.AttachmentReference {
				return vk.AttachmentReference{
					Attachment: uint32(idx),
					Layout:     vk.ImageLayoutColorAttachmentOptimal,
				}
			}),
			PDepthStencilAttachment: depthRef,
		})
		subpassIndices[subpass.Name] = idx
	}

	dependencies := make([]vk.SubpassDependency, len(args.Dependencies))
	for idx, dependency := range args.Dependencies {
		dependencies[idx] = vk.SubpassDependency{
			SrcSubpass:      uint32(subpassIndices[dependency.Src]),
			DstSubpass:      uint32(subpassIndices[dependency.Dst]),
			SrcStageMask:    dependency.SrcStageMask,
			SrcAccessMask:   dependency.SrcAccessMask,
			DstStageMask:    dependency.DstStageMask,
			DstAccessMask:   dependency.SrcAccessMask,
			DependencyFlags: dependency.Flags,
		}
	}

	info := vk.RenderPassCreateInfo{
		SType:           vk.StructureTypeRenderPassCreateInfo,
		AttachmentCount: uint32(len(descriptions)),
		PAttachments:    descriptions,
		SubpassCount:    uint32(len(subpasses)),
		PSubpasses:      subpasses,
		DependencyCount: uint32(len(dependencies)),
		PDependencies:   dependencies,
	}

	var ptr vk.RenderPass
	vk.CreateRenderPass(device.Ptr(), &info, nil, &ptr)

	framebuffers := make([]framebuffer.T, args.Frames)
	for i := range framebuffers {
		fbviews := make([]image.View, 0, len(descriptions))
		for _, name := range attachmentOrder {
			attachment := attachments[name]
			fbviews = append(fbviews, attachment.View(i))
		}
		if depth != nil {
			fbviews = append(fbviews, depth.View(i))
		}
		framebuffers[i] = framebuffer.New(device, args.Width, args.Height, ptr, fbviews)
	}

	return &renderpass{
		device:       device,
		ptr:          ptr,
		framebuffers: framebuffers,
		depth:        depth,
		order:        attachmentOrder,
		attachments:  attachments,
		clear:        clear,
	}
}

func (r *renderpass) Ptr() vk.RenderPass { return r.ptr }
func (r *renderpass) Depth() Attachment  { return r.depth }

func (r *renderpass) Attachment(name string) Attachment {
	return r.attachments[name]
}

func (r *renderpass) Framebuffer(frame int) framebuffer.T {
	return r.framebuffers[frame]
}

func (r *renderpass) Clear() []vk.ClearValue {
	return r.clear
}

func (r *renderpass) Attachments() []Attachment {
	return util.Map(r.order, func(i int, name string) Attachment { return r.attachments[name] })
}

func (r *renderpass) Destroy() {
	if r.ptr != nil {
		vk.DestroyRenderPass(r.device.Ptr(), r.ptr, nil)
		r.ptr = nil
	}

	for _, fb := range r.framebuffers {
		fb.Destroy()
	}

	if r.depth != nil {
		r.depth.Destroy()
	}

	for _, attachment := range r.attachments {
		attachment.Destroy()
	}
}