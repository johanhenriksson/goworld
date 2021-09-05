package engine

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type DrawPass interface {
	Draw(*Scene)
	Resize(int, int)
}

// DrawArgs holds the arguments used to perform a draw pass.
// Includes the various transformation matrices and position of the camera.
type DrawArgs struct {
	VP         mat4.T
	MVP        mat4.T
	Projection mat4.T
	View       mat4.T
	Transform  mat4.T
	Position   vec3.T
	Pass       render.Pass
}

// Apply the effects of a transform
func (d DrawArgs) Apply(t mat4.T) DrawArgs {
	d.Transform = d.Transform.Mul(&t)
	d.MVP = d.VP.Mul(&d.Transform)
	return d
}

func (d DrawArgs) Set(t transform.T) DrawArgs {
	d.Transform = t.World()
	d.MVP = d.VP.Mul(&d.Transform)
	return d
}

func ArgsFromCamera(cam camera.T) DrawArgs {
	return DrawArgs{
		Projection: cam.Projection(),
		View:       cam.View(),
		VP:         cam.ViewProj(),
		MVP:        cam.ViewProj(),
		Transform:  mat4.Ident(),
		Position:   cam.Transform().WorldPosition(),
	}
}
