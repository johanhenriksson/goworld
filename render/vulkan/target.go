package vulkan

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/swapchain"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// the render target interfaces & implementations probably dont belong in this package long-term

type TargetSize struct {
	Width  int
	Height int
	Frames int
	Scale  float32
}

type Target interface {
	Size() TargetSize
	Scale() float32
	Width() int
	Height() int
	Frames() int

	Surfaces() []image.T
	SurfaceFormat() core1_0.Format
	Aquire() (*swapchain.Context, error)
	Present(*swapchain.Context)

	Destroy()
}

type renderTarget struct {
	size     TargetSize
	format   core1_0.Format
	usage    core1_0.ImageUsageFlags
	surfaces []image.T
	context  *swapchain.Context
}

func NewDepthTarget(device device.T, key string, size TargetSize) Target {
	format := device.GetDepthFormat()
	usage := core1_0.ImageUsageSampled | core1_0.ImageUsageDepthStencilAttachment | core1_0.ImageUsageInputAttachment
	target, err := NewRenderTarget(device, key, format, usage, size)
	if err != nil {
		panic(err)
	}
	return target
}

func NewColorTarget(device device.T, key string, format core1_0.Format, size TargetSize) Target {
	usage := core1_0.ImageUsageSampled | core1_0.ImageUsageColorAttachment | core1_0.ImageUsageInputAttachment | core1_0.ImageUsageTransferSrc
	target, err := NewRenderTarget(device, key, format, usage, size)
	if err != nil {
		panic(err)
	}
	return target
}

func NewRenderTarget(device device.T, key string, format core1_0.Format, usage core1_0.ImageUsageFlags, size TargetSize) (Target, error) {
	var err error
	outputs := make([]image.T, size.Frames)
	for i := 0; i < size.Frames; i++ {
		outputs[i], err = image.New2D(device, fmt.Sprintf("%s:%d", key, i), size.Width, size.Height, format, false, usage)
		if err != nil {
			return nil, err
		}
	}

	return &renderTarget{
		size:     size,
		format:   format,
		usage:    usage,
		surfaces: outputs,
		context:  swapchain.DummyContext(),
	}, nil
}

func (r *renderTarget) Frames() int    { return len(r.surfaces) }
func (r *renderTarget) Width() int     { return r.size.Width }
func (r *renderTarget) Height() int    { return r.size.Height }
func (r *renderTarget) Scale() float32 { return r.size.Scale }

func (r *renderTarget) Size() TargetSize {
	return r.size
}

func (r *renderTarget) Destroy() {
	for _, output := range r.surfaces {
		output.Destroy()
	}
	r.surfaces = nil
}

func (i *renderTarget) Surfaces() []image.T           { return i.surfaces }
func (i *renderTarget) SurfaceFormat() core1_0.Format { return i.format }

func (i *renderTarget) Aquire() (*swapchain.Context, error) {
	return i.context, nil
}

func (b *renderTarget) Present(*swapchain.Context) {

}
