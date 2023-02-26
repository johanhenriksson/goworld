package renderpass

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/johanhenriksson/goworld/util"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/core/v2/driver"
)

type T interface {
	device.Resource[core1_0.RenderPass]

	Depth() attachment.T
	Attachment(name attachment.Name) attachment.T
	Attachments() []attachment.T
	Subpass(name Name) Subpass
	Clear() []core1_0.ClearValue
}

type renderpass struct {
	device      device.T
	ptr         core1_0.RenderPass
	subpasses   []Subpass
	passIndices map[Name]int
	attachments []attachment.T
	depth       attachment.T
	indices     map[attachment.Name]int
	clear       []core1_0.ClearValue
}

func New(device device.T, args Args) T {
	clear := make([]core1_0.ClearValue, 0, len(args.ColorAttachments)+1)
	attachments := make([]attachment.T, len(args.ColorAttachments))
	attachmentIndices := make(map[attachment.Name]int)

	log.Println("create renderpass", args.Name)
	log.Println("attachments")
	for index, desc := range args.ColorAttachments {
		attachment := attachment.NewColor(device, desc)
		clear = append(clear, attachment.Clear())
		attachments[index] = attachment
		attachmentIndices[attachment.Name()] = index
		log.Printf("  %d: %s", index, desc.Name)
	}

	var depth attachment.T
	if args.DepthAttachment != nil {
		index := len(attachments)
		attachmentIndices[attachment.DepthName] = index
		depth = attachment.NewDepth(device, *args.DepthAttachment)
		clear = append(clear, depth.Clear())
		log.Printf("  %d: %s", index, attachment.DepthName)
	}

	descriptions := make([]core1_0.AttachmentDescription, 0, len(args.ColorAttachments)+1)
	for _, attachment := range attachments {
		descriptions = append(descriptions, attachment.Description())
	}
	if depth != nil {
		descriptions = append(descriptions, depth.Description())
	}

	subpasses := make([]core1_0.SubpassDescription, 0, len(args.Subpasses))
	subpassIndices := make(map[Name]int)

	for idx, subpass := range args.Subpasses {
		log.Println("subpass", idx)

		var depthRef *core1_0.AttachmentReference
		if depth != nil && subpass.Depth {
			idx := attachmentIndices[attachment.DepthName]
			depthRef = &core1_0.AttachmentReference{
				Attachment: idx,
				Layout:     core1_0.ImageLayoutDepthStencilAttachmentOptimal,
			}
			log.Printf("  depth -> %s (%d)\n", attachment.DepthName, idx)
		}

		subpasses = append(subpasses, core1_0.SubpassDescription{
			PipelineBindPoint: core1_0.PipelineBindPointGraphics,

			ColorAttachments: util.MapIdx(
				subpass.ColorAttachments,
				func(name attachment.Name, i int) core1_0.AttachmentReference {
					idx := attachmentIndices[name]
					log.Printf("  color %d -> %s (%d)\n", i, name, idx)
					return core1_0.AttachmentReference{
						Attachment: idx,
						Layout:     core1_0.ImageLayoutColorAttachmentOptimal,
					}
				}),

			InputAttachments: util.MapIdx(
				subpass.InputAttachments,
				func(name attachment.Name, i int) core1_0.AttachmentReference {
					idx := attachmentIndices[name]
					log.Printf("  input %d -> %s (%d)\n", i, name, idx)
					return core1_0.AttachmentReference{
						Attachment: idx,
						Layout:     core1_0.ImageLayoutShaderReadOnlyOptimal,
					}
				}),

			DepthStencilAttachment: depthRef,
		})

		subpassIndices[subpass.Name] = idx
		args.Subpasses[idx].index = idx
	}

	dependencies := make([]core1_0.SubpassDependency, len(args.Dependencies))
	for idx, dependency := range args.Dependencies {
		src := core1_0.SubpassExternal
		if dependency.Src != ExternalSubpass {
			src = subpassIndices[dependency.Src]
		}
		dst := core1_0.SubpassExternal
		if dependency.Dst != ExternalSubpass {
			dst = subpassIndices[dependency.Dst]
		}
		dependencies[idx] = core1_0.SubpassDependency{
			SrcSubpass:      src,
			DstSubpass:      dst,
			SrcStageMask:    core1_0.PipelineStageFlags(dependency.SrcStageMask),
			SrcAccessMask:   core1_0.AccessFlags(dependency.SrcAccessMask),
			DstStageMask:    core1_0.PipelineStageFlags(dependency.DstStageMask),
			DstAccessMask:   core1_0.AccessFlags(dependency.DstAccessMask),
			DependencyFlags: core1_0.DependencyFlags(dependency.Flags),
		}
	}

	ptr, _, err := device.Ptr().CreateRenderPass(nil, core1_0.RenderPassCreateInfo{
		Attachments:         descriptions,
		Subpasses:           subpasses,
		SubpassDependencies: dependencies,
	})
	if err != nil {
		panic(err)
	}

	// set object name
	device.SetDebugObjectName(driver.VulkanHandle(ptr.Handle()), core1_0.ObjectTypeRenderPass, args.Name)

	return &renderpass{
		device:      device,
		ptr:         ptr,
		depth:       depth,
		indices:     attachmentIndices,
		attachments: attachments,
		passIndices: subpassIndices,
		subpasses:   args.Subpasses,
		clear:       clear,
	}
}

func (r *renderpass) Ptr() core1_0.RenderPass { return r.ptr }
func (r *renderpass) Depth() attachment.T     { return r.depth }

func (r *renderpass) Attachment(name attachment.Name) attachment.T {
	if name == attachment.DepthName {
		return r.depth
	}
	index := r.indices[name]
	return r.attachments[index]
}

func (r *renderpass) Clear() []core1_0.ClearValue {
	return r.clear
}

func (r *renderpass) Attachments() []attachment.T {
	return r.attachments
}

func (r *renderpass) Subpass(name Name) Subpass {
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
		r.ptr.Destroy(nil)
		r.ptr = nil
	}
}
