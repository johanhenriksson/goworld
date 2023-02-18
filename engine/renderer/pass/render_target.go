package pass

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/renderpass/attachment"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type RenderTarget interface {
	Destroy()
	Width() int
	Height() int
	Output() image.T
	Depth() image.T
}

// renderTarget holds color and/or depth textures to render to.
type renderTarget struct {
	width  int
	height int
	output image.T
	depth  image.T
}

func NewRenderTarget(device device.T, width, height int, outputFormat, depthFormat core1_0.Format) (RenderTarget, error) {
	output, err := image.New2D(device, "output", width, height, outputFormat,
		core1_0.ImageUsageSampled|core1_0.ImageUsageColorAttachment|core1_0.ImageUsageInputAttachment)
	if err != nil {
		return nil, err
	}

	depth, err := image.New2D(device, "depth", width, height, depthFormat,
		core1_0.ImageUsageSampled|core1_0.ImageUsageDepthStencilAttachment|core1_0.ImageUsageInputAttachment)
	if err != nil {
		return nil, err
	}

	return &renderTarget{
		width:  width,
		height: height,
		output: output,
		depth:  depth,
	}, nil
}

func (r *renderTarget) Width() int  { return r.width }
func (r *renderTarget) Height() int { return r.height }

func (r *renderTarget) OutputAttachment() attachment.Color {
	return attachment.Color{
		Name:        OutputAttachment,
		Format:      r.output.Format(),
		LoadOp:      core1_0.AttachmentLoadOpLoad,
		StoreOp:     core1_0.AttachmentStoreOpStore,
		FinalLayout: core1_0.ImageLayoutShaderReadOnlyOptimal,
		Usage:       core1_0.ImageUsageSampled,
		Blend:       attachment.BlendMultiply,

		Allocator: attachment.FromImageArray([]image.T{
			r.output,
		}),
	}
}

func (r *renderTarget) Output() image.T {
	return r.output
}

func (r *renderTarget) DepthAttachment() *attachment.Depth {
	return &attachment.Depth{
		LoadOp:        core1_0.AttachmentLoadOpLoad,
		StencilLoadOp: core1_0.AttachmentLoadOpLoad,
		StoreOp:       core1_0.AttachmentStoreOpStore,
		FinalLayout:   core1_0.ImageLayoutShaderReadOnlyOptimal,
		Usage:         core1_0.ImageUsageInputAttachment,
		Format:        r.depth.Format(),

		Allocator: attachment.FromImageArray([]image.T{
			r.depth,
		}),
	}
}

func (r *renderTarget) Depth() image.T {
	return r.depth
}

func (r *renderTarget) Destroy() {
	r.output.Destroy()
	r.output = nil
	r.depth.Destroy()
	r.depth = nil
}
