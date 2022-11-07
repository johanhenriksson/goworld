package widget

import (
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
)

type Renderer[W T] interface {
	Draw(DrawArgs, W)
}

type DrawArgs struct {
	Commands  command.Recorder
	Meshes    cache.MeshCache
	Textures  cache.SamplerCache
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
