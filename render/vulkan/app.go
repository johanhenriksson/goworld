package vulkan

import (
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/swapchain"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Target interface {
	Scale() float32
	Width() int
	Height() int
	Frames() int

	Surfaces() []image.T
	SurfaceFormat() core1_0.Format
	Aquire() (*swapchain.Context, error)
	Present(command.Worker, *swapchain.Context)
}
