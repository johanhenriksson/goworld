package engine

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
)

type DrawPass interface {
	Draw(render.Args, object.T)
}

func CreateRenderArgs(window window.T, cam camera.T) render.Args {
	w, h := window.Size()

	return render.Args{
		Projection: cam.Projection(),
		View:       cam.View(),
		VP:         cam.ViewProj(),
		MVP:        cam.ViewProj(),
		Transform:  mat4.Ident(),
		Position:   cam.Transform().WorldPosition(),
		Clear:      cam.ClearColor(),

		Viewport: render.Screen{
			Width:  w,
			Height: h,
			Scale:  window.Scale(),
		},
	}
}
