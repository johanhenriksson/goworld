package render

import (
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Args holds the arguments used to perform a draw pass.
// Includes the various transformation matrices and position of the camera.
type Args struct {
	VP         mat4.T
	MVP        mat4.T
	Projection mat4.T
	View       mat4.T
	Transform  mat4.T
	Position   vec3.T
	Viewport   Screen
}

type Screen struct {
	Width       int
	Height      int
	FrameWidth  int
	FrameHeight int
	Scale       float32
}

// Apply the effects of a transform
func (d Args) Apply(t mat4.T) Args {
	d.Transform = d.Transform.Mul(&t)
	d.MVP = d.VP.Mul(&d.Transform)
	return d
}

func (d Args) Set(t transform.T) Args {
	d.Transform = t.World()
	d.MVP = d.VP.Mul(&d.Transform)
	return d
}
