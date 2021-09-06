package engine

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/render"
)

type DrawPass interface {
	Draw(scene.T)
	Resize(int, int)
}

func ArgsFromCamera(cam camera.T) render.Args {
	return render.Args{
		Projection: cam.Projection(),
		View:       cam.View(),
		VP:         cam.ViewProj(),
		MVP:        cam.ViewProj(),
		Transform:  mat4.Ident(),
		Position:   cam.Transform().WorldPosition(),
	}
}
