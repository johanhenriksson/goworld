package renderpass

import (
	"fmt"
	"log"

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
	Subpass(name string) Subpass
	Clear() []vk.ClearValue
}

type renderpass struct {
	device       device.T
	ptr          vk.RenderPass
	framebuffers []framebuffer.T
	subpasses    []Subpass
	passIndices  map[string]int
	attachments  []Attachment
	depth        Attachment
	indices      map[string]int
	clear        []vk.ClearValue
}

func New(device device.T, args Args) T {
	clear := make([]vk.ClearValue, 0, len(args.ColorAttachments)+1)
	attachments := make([]Attachment, len(args.ColorAttachments))
	attachmentIndices := make(map[string]int)

	log.Println("attachments")
	for index, desc := range args.ColorAttachments {
		attachment := NewColorAttachment(device, desc, args.Frames, args.Width, args.Height)
		clear = append(clear, attachment.Clear())
		attachments[index] = attachment
		attachmentIndices[attachment.Name()] = index
		log.Printf("  %d: %s", index, desc.Name)
	}

	var depth Attachment
	if args.DepthAttachment != nil {
		index := len(attachments)
		attachmentIndices["depth"] = index
		depth = NewDepthAttachment(device, *args.DepthAttachment, args.Frames, args.Width, args.Height, index)
		clear = append(clear, depth.Clear())
		log.Printf("  %d: %s", index, "depth")
	}

	descriptions := make([]vk.AttachmentDescription, 0, len(args.ColorAttachments)+1)
	for _, attachment := range attachments {
		descriptions = append(descriptions, attachment.Description())
	}
	if depth != nil {
		descriptions = append(descriptions, depth.Description())
	}

	subpasses := make([]vk.SubpassDescription, 0, len(args.Subpasses))
	subpassIndices := make(map[string]int)

	for idx, subpass := range args.Subpasses {
		log.Println("subpass", idx)

		var depthRef *vk.AttachmentReference
		if depth != nil && subpass.Depth {
			idx := attachmentIndices["depth"]
			depthRef = &vk.AttachmentReference{
				Attachment: uint32(idx),
				Layout:     vk.ImageLayoutDepthStencilAttachmentOptimal,
			}
			log.Printf("  depth -> %s (%d)\n", "depth", idx)
		}

		subpasses = append(subpasses, vk.SubpassDescription{
			PipelineBindPoint: vk.PipelineBindPointGraphics,

			ColorAttachmentCount: uint32(len(subpass.ColorAttachments)),
			PColorAttachments: util.MapIdx(
				subpass.ColorAttachments,
				func(name string, i int) vk.AttachmentReference {
					idx := attachmentIndices[name]
					log.Printf("  color %d -> %s (%d)\n", i, name, idx)
					return vk.AttachmentReference{
						Attachment: uint32(idx),
						Layout:     vk.ImageLayoutColorAttachmentOptimal,
					}
				}),

			InputAttachmentCount: uint32(len(subpass.InputAttachments)),
			PInputAttachments: util.MapIdx(
				subpass.InputAttachments,
				func(name string, i int) vk.AttachmentReference {
					idx := attachmentIndices[name]
					log.Printf("  input %d -> %s (%d)\n", i, name, idx)
					return vk.AttachmentReference{
						Attachment: uint32(idx),
						Layout:     vk.ImageLayoutShaderReadOnlyOptimal,
					}
				}),

			PDepthStencilAttachment: depthRef,
		})

		subpassIndices[subpass.Name] = idx
		args.Subpasses[idx].index = idx
	}

	// subpassIndices["external"] = int(vk.SubpassExternal)

	dependencies := make([]vk.SubpassDependency, len(args.Dependencies))
	for idx, dependency := range args.Dependencies {
		src := vk.SubpassExternal
		if dependency.Src != "external" {
			src = uint32(subpassIndices[dependency.Src])
		}
		dst := vk.SubpassExternal
		if dependency.Dst != "external" {
			dst = uint32(subpassIndices[dependency.Dst])
		}
		dependencies[idx] = vk.SubpassDependency{
			SrcSubpass:      src,
			DstSubpass:      dst,
			SrcStageMask:    vk.PipelineStageFlags(dependency.SrcStageMask),
			SrcAccessMask:   vk.AccessFlags(dependency.SrcAccessMask),
			DstStageMask:    vk.PipelineStageFlags(dependency.DstStageMask),
			DstAccessMask:   vk.AccessFlags(dependency.DstAccessMask),
			DependencyFlags: vk.DependencyFlags(dependency.Flags),
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
		fbviews := util.Map(attachments, func(attachment Attachment) image.View { return attachment.View(i) })
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
		indices:      attachmentIndices,
		attachments:  attachments,
		passIndices:  subpassIndices,
		subpasses:    args.Subpasses,
		clear:        clear,
	}
}

func (r *renderpass) Ptr() vk.RenderPass { return r.ptr }
func (r *renderpass) Depth() Attachment  { return r.depth }

func (r *renderpass) Attachment(name string) Attachment {
	if name == "depth" {
		return r.depth
	}
	index := r.indices[name]
	return r.attachments[index]
}

func (r *renderpass) Framebuffer(frame int) framebuffer.T {
	return r.framebuffers[frame%len(r.framebuffers)]
}

func (r *renderpass) Clear() []vk.ClearValue {
	return r.clear
}

func (r *renderpass) Attachments() []Attachment {
	return r.attachments
}

func (r *renderpass) Subpass(name string) Subpass {
	if name == "" {
		return r.subpasses[0]
	}
	idx, exists := r.passIndices[name]
	if !exists {
		panic(fmt.Sprintf("unknown subpass %s", name))
	}
	return r.subpasses[idx]
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
