package engine

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Transform represents a 3D transformation
type Transform struct {
	Matrix   mat4.T
	Position vec3.T
	Rotation vec3.T
	Scale    vec3.T
	Forward  vec3.T
	Right    vec3.T
	Up       vec3.T
}

// NewTransform creates a new 3D transform
func NewTransform(position, rotation, scale vec3.T) *Transform {
	t := &Transform{
		Matrix:   mat4.Ident(),
		Position: position,
		Rotation: rotation,
		Scale:    scale,
	}
	t.Update(0)
	return t
}

func Identity() *Transform {
	return NewTransform(vec3.Zero, vec3.Zero, vec3.One)
}

// Update transform matrix and its right/up/forward vectors
func (t *Transform) Update(dt float32) {
	// Update transform
	m := mat4.Transform(t.Position, t.Rotation, t.Scale)

	// Grab axis vectors from transformation matrix
	t.Up = m.Up()
	t.Right = m.Right()
	t.Forward = m.Forward()

	// Update transformation matrix
	t.Matrix = m
}

// Translate this transform by the given offset
func (t *Transform) Translate(offset vec3.T) {
	t.Position = t.Position.Add(offset)
}

// TransformPoint transforms a point into this coordinate system
func (t *Transform) TransformPoint(point vec3.T) vec3.T {
	return t.Matrix.TransformPoint(point)
}

// TransformDir transforms a direction vector into this coordinate system
func (t *Transform) TransformDir(dir vec3.T) vec3.T {
	return t.Matrix.TransformDir(dir)
}
