package vulkan

import (
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/swapchain"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type imageTarget struct {
	T
	image   image.T
	context *swapchain.Context
}

func NewImageTarget(backend T, img image.T) App {
	return &imageTarget{
		T:     backend,
		image: img,
	}
}

func (i *imageTarget) Frames() int                   { return 1 }
func (i *imageTarget) Scale() float32                { return 1 }
func (i *imageTarget) Width() int                    { return i.image.Width() }
func (i *imageTarget) Height() int                   { return i.image.Height() }
func (i *imageTarget) Surfaces() []image.T           { return []image.T{i.image} }
func (i *imageTarget) SurfaceFormat() core1_0.Format { return i.image.Format() }

func (i *imageTarget) Aquire() (*swapchain.Context, error) {
	return i.context, nil
}

func (b *imageTarget) Present(command.Worker, *swapchain.Context) {}

func (b *imageTarget) Destroy() {
	b.image.Destroy()
}
