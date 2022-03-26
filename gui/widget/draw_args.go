package widget

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/cache"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
)

type Renderer[W T] interface {
	Draw(DrawArgs, W)
	Destroy()
}

type DrawArgs struct {
	Commands  command.Recorder
	Meshes    cache.MeshCache
	Textures  cache.TextureCache
	ViewProj  mat4.T
	Transform mat4.T
	Position  vec3.T
	Viewport  render.Screen
}

type Constants struct {
	Viewport mat4.T
	Model    mat4.T
	Texture  int
}
