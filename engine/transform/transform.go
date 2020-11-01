package transform

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Transform represents a 3D transformation
type T struct {
	World    mat4.T
	Local    mat4.T
	Position vec3.T
	Rotation vec3.T
	Scale    vec3.T
	Forward  vec3.T
	Right    vec3.T
	Up       vec3.T
}

// NewTransform creates a new 3D transform
func New(position, rotation, scale vec3.T) *T {
	t := &T{
		World:    mat4.Ident(),
		Local:    mat4.Ident(),
		Position: position,
		Rotation: rotation,
		Scale:    scale,
	}
	t.Update(&t.World)
	return t
}

// Identity returns a new transform that does nothing.
func Identity() *T {
	return New(vec3.Zero, vec3.Zero, vec3.One)
}

// Update transform matrix and its right/up/forward vectors
func (t *T) Update(parent *mat4.T) {
	// Update transform
	m := mat4.Transform(t.Position, t.Rotation, t.Scale)

	// Update parent -> local transformation matrix
	t.Local = m

	// Update local -> world matrix
	if parent != nil {
		t.World = m.Mul(parent)
	} else {
		t.World = m
	}

	// Grab axis vectors from transformation matrix
	t.Up = t.World.Up()
	t.Right = t.World.Right()
	t.Forward = t.World.Forward()
}

// Translate this transform by the given offset
func (t *T) Translate(offset vec3.T) {
	t.Position = t.Position.Add(offset)
}

// TransformPoint transforms a world point into this coordinate system
func (t *T) TransformPoint(point vec3.T) vec3.T {
	return t.World.TransformPoint(point)
}

// TransformDir transforms a world direction vector into this coordinate system
func (t *T) TransformDir(dir vec3.T) vec3.T {
	return t.World.TransformDir(dir)
}
