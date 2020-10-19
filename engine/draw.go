package engine

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// DrawPass is a type identifier for a render draw pass.
type DrawPass int32

// Various draw pass constants
const (
	DrawGeometry DrawPass = iota
	DrawShadow
	DrawLines
	DrawParticles
	DrawForward
)

// DrawArgs holds the arguments used to perform a draw pass.
// Includes the various transformation matrices and position of the camera.
type DrawArgs struct {
	VP         mat4.T
	MVP        mat4.T
	Projection mat4.T
	View       mat4.T
	Transform  mat4.T
	Position   vec3.T
	Pass       DrawPass
}

// Apply the effects of a transform
func (d DrawArgs) Apply(t *Transform) DrawArgs {
	d.Transform = d.Transform.Mul(&t.Matrix)
	d.MVP = d.VP.Mul(&d.Transform)
	return d
}
