package pass

import (
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/swapchain"
	"github.com/vkngwrapper/core/v2/core1_0"
)

type RenderTarget interface {
	Destroy()
	Frames() int
	Width() int
	Height() int
	Output() []image.T
	Depth() []image.T

	Scale() float32
	Surfaces() []image.T
	SurfaceFormat() core1_0.Format
	Aquire() (*swapchain.Context, error)
	Present(command.Worker, *swapchain.Context)
}

// renderTarget holds color and/or depth textures to render to.
type renderTarget struct {
	width  int
	height int
	output []image.T
	depth  []image.T
}

func NewRenderTarget(device device.T, width, height, frames int, outputFormat, depthFormat core1_0.Format) (RenderTarget, error) {
	var err error
	outputs := make([]image.T, frames)
	for i := 0; i < frames; i++ {
		outputs[i], err = image.New2D(device, "output", width, height, outputFormat,
			core1_0.ImageUsageSampled|core1_0.ImageUsageColorAttachment|core1_0.ImageUsageInputAttachment)
		if err != nil {
			return nil, err
		}
	}

	var depths []image.T
	if depthFormat != core1_0.FormatUndefined {
		depths = make([]image.T, frames)
		for i := 0; i < frames; i++ {
			depths[i], err = image.New2D(device, "depth", width, height, depthFormat,
				core1_0.ImageUsageSampled|core1_0.ImageUsageDepthStencilAttachment|core1_0.ImageUsageInputAttachment)
			if err != nil {
				return nil, err
			}
		}
	}

	return &renderTarget{
		width:  width,
		height: height,
		output: outputs,
		depth:  depths,
	}, nil
}

func (r *renderTarget) Frames() int    { return len(r.output) }
func (r *renderTarget) Width() int     { return r.width }
func (r *renderTarget) Height() int    { return r.height }
func (r *renderTarget) Scale() float32 { return 1 }

func (r *renderTarget) Output() []image.T {
	return r.output
}

func (r *renderTarget) Depth() []image.T {
	return r.depth
}

func (r *renderTarget) Destroy() {
	for _, output := range r.output {
		output.Destroy()
	}
	r.output = nil

	for _, depth := range r.depth {
		depth.Destroy()
	}
	r.depth = nil
}

func (i *renderTarget) Surfaces() []image.T           { return i.output }
func (i *renderTarget) SurfaceFormat() core1_0.Format { return i.output[0].Format() }

func (i *renderTarget) Aquire() (*swapchain.Context, error) {
	return nil, nil
}

func (b *renderTarget) Present(command.Worker, *swapchain.Context) {}
