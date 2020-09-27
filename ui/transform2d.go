package ui

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Transform2D struct {
	Matrix   mat4.T
	Position vec3.T
	Scale    vec2.T
	Rotation float32
}

/* Creates a new 2D transform */
func CreateTransform2D(x, y, z float32) *Transform2D {
	t := &Transform2D{
		Matrix:   mat4.Ident(),
		Position: vec3.T{x, y, z},
		Scale:    vec2.One,
		Rotation: 0.0,
	}
	t.Update(0)
	return t
}

func (t *Transform2D) Update(dt float32) {
	t.Matrix = mat4.Transform(t.Position, vec3.T{0, 0, t.Rotation}, vec3.Extend(t.Scale, 1))
}
