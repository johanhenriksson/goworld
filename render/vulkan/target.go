package vulkan

import (
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
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

	Worker(int) command.Worker
	Transferer() command.Worker

	Meshes() cache.MeshCache
	Textures() cache.TextureCache
}