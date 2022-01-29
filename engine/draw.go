package engine

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
)

type DrawPass interface {
	Draw(render.Args, scene.T)
}

func ArgsWithCamera(args render.Args, cam camera.T) render.Args {
	args.Projection = cam.Projection()
	args.View = cam.View()
	args.VP = cam.ViewProj()
	args.MVP = cam.ViewProj()
	args.Transform = mat4.Ident()
	args.Position = cam.Transform().WorldPosition()
	return args
}

func ArgsFromWindow(window window.T) render.Args {
	w, h := window.Size()
	fw, fh := window.BufferSize()
	return render.Args{
		Viewport: render.Viewport{
			Width:       w,
			Height:      h,
			FrameWidth:  fw,
			FrameHeight: fh,
			Scale:       window.Scale(),
		},
	}
}
