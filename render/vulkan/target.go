package vulkan

import (
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/swapchain"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type Target interface {
	Device() device.T
	Destroy()

	Scale() float32
	Width() int
	Height() int
	Frames() int

	Surfaces() []image.T
	SurfaceFormat() core1_0.Format
	Aquire() (swapchain.Context, error)

	Worker(int) command.Worker
	Transferer() command.Worker

	Pool() descriptor.Pool
	Meshes() cache.MeshCache
	Textures() cache.TextureCache
}
