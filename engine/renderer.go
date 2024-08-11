package engine

import (
	"image"

	"github.com/johanhenriksson/goworld/core/object"
)

type RendererFunc func(App, Target) Renderer

type Renderer interface {
	Draw(scene object.Object, time, delta float32)
	Recreate()
	Screengrab() *image.RGBA
	Destroy()
}
