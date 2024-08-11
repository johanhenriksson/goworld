package engine

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/instance"
)

type SceneFunc func(object.Object)

type Args struct {
	Title  string
	Width  int
	Height int
}

type App interface {
	Instance() instance.T
	Device() device.T
	Destroy()

	Worker() command.Worker
	Flush()

	Pool() descriptor.Pool
	Meshes() cache.MeshCache
	Textures() cache.TextureCache
	Shaders() cache.ShaderCache
}
