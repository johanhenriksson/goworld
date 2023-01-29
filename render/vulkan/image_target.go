package vulkan

import (
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/swapchain"
	vk "github.com/vulkan-go/vulkan"
)

type imageTarget struct {
	T
	image image.T
}

func NewImageTarget(backend T, img image.T) Target {
	return &imageTarget{
		T:     backend,
		image: img,
	}
}

func (i *imageTarget) Frames() int              { return 1 }
func (i *imageTarget) Scale() float32           { return 1 }
func (i *imageTarget) Width() int               { return i.image.Width() }
func (i *imageTarget) Height() int              { return i.image.Height() }
func (i *imageTarget) Surfaces() []image.T      { return []image.T{i.image} }
func (i *imageTarget) SurfaceFormat() vk.Format { return i.image.Format() }

func (i *imageTarget) Aquire() (swapchain.Context, error) {
	return swapchain.Context{}, nil
}

func (b *imageTarget) Present() {
}

func (b *imageTarget) Destroy() {
	b.image.Destroy()
}
