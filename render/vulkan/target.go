package vulkan

import (
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/swapchain"
	vk "github.com/vulkan-go/vulkan"
)

type Target interface {
	Device() device.T
	Destroy()

	Scale() float32
	Width() int
	Height() int
	Frames() int

	Surfaces() []image.T
	SurfaceFormat() vk.Format
	Aquire() (swapchain.Context, error)

	Worker(int) command.Worker
	Transferer() command.Worker

	Pool() descriptor.Pool
	Meshes() cache.MeshCache
	Textures() cache.TextureCache
}
