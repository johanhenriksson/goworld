package engine

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type DrawPass int32

const (
	DrawGeometry DrawPass = iota
	DrawShadow
	DrawLines
	DrawParticles
	DrawForward
)

/** UI Component render interface */
type Drawable interface {
	Draw(DrawArgs)
}

/** Passed to Drawables on render */
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
