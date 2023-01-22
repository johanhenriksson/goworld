package vulkan

import (
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/swapchain"
)

type Target interface {
	Device() device.T

	Scale() float32
	Width() int
	Height() int
	Frames() int

	Swapchain() swapchain.T
	Aquire() (swapchain.Context, error)
	Present()
	Worker(int) command.Worker
	Transferer() command.Worker

	Meshes() cache.MeshCache
	Textures() cache.TextureCache
}
