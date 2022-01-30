package engine

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
)

type DrawPass interface {
	// Prepare()
	Draw(render.Args, scene.T)
}

func CreateRenderArgs(window window.T, cam camera.T) render.Args {
	w, h := window.Size()
	fw, fh := window.BufferSize()

	return render.Args{
		Projection: cam.Projection(),
		View:       cam.View(),
		VP:         cam.ViewProj(),
		MVP:        cam.ViewProj(),
		Transform:  mat4.Ident(),
		Position:   cam.Transform().WorldPosition(),

		Viewport: render.Viewport{
			Width:       w,
			Height:      h,
			FrameWidth:  fw,
			FrameHeight: fh,
			Scale:       window.Scale(),
		},
	}
}
