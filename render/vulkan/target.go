package vulkan

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/swapchain"

	"github.com/vkngwrapper/core/v2/core1_0"
)

// the render target interfaces & implementations probably dont belong in this package long-term

type Target interface {
	Scale() float32
	Width() int
	Height() int
	Frames() int

	Surfaces() []image.T
	SurfaceFormat() core1_0.Format
	Aquire() (*swapchain.Context, error)
	Present(command.Worker, *swapchain.Context)

	Destroy()
}

type renderTarget struct {
	width   int
	height  int
	scale   float32
	format  core1_0.Format
	usage   core1_0.ImageUsageFlags
	surface []image.T
	context *swapchain.Context
}

func NewDepthTarget(device device.T, key string, width, height, frames int, scale float32) (Target, error) {
	format := device.GetDepthFormat()
	usage := core1_0.ImageUsageSampled | core1_0.ImageUsageDepthStencilAttachment | core1_0.ImageUsageInputAttachment
	return NewRenderTarget(device, key, width, height, frames, scale, format, usage)
}

func NewColorTarget(device device.T, key string, width, height, frames int, scale float32, format core1_0.Format) (Target, error) {
	usage := core1_0.ImageUsageSampled | core1_0.ImageUsageColorAttachment | core1_0.ImageUsageInputAttachment | core1_0.ImageUsageTransferSrc
	return NewRenderTarget(device, key, width, height, frames, scale, format, usage)
}

func NewRenderTarget(device device.T, key string, width, height, frames int, scale float32, format core1_0.Format, usage core1_0.ImageUsageFlags) (Target, error) {
	var err error
	outputs := make([]image.T, frames)
	for i := 0; i < frames; i++ {
		outputs[i], err = image.New2D(device, fmt.Sprintf("%s:%d", key, i), width, height, format, usage)
		if err != nil {
			return nil, err
		}
	}

	return &renderTarget{
		width:   width,
		height:  height,
		scale:   scale,
		format:  format,
		usage:   usage,
		surface: outputs,
		context: swapchain.DummyContext(),
	}, nil
}

func (r *renderTarget) Frames() int    { return len(r.surface) }
func (r *renderTarget) Width() int     { return r.width }
func (r *renderTarget) Height() int    { return r.height }
func (r *renderTarget) Scale() float32 { return r.scale }

func (r *renderTarget) Destroy() {
	for _, output := range r.surface {
		output.Destroy()
	}
	r.surface = nil
}

func (i *renderTarget) Surfaces() []image.T           { return i.surface }
func (i *renderTarget) SurfaceFormat() core1_0.Format { return i.format }

func (i *renderTarget) Aquire() (*swapchain.Context, error) {
	i.context.Aquire()
	return i.context, nil
}

func (b *renderTarget) Present(command.Worker, *swapchain.Context) {

}
