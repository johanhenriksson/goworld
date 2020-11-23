package transform

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Transform represents a 3D transformation
type T struct {
	World mat4.T
	Local mat4.T

	forward  vec3.T
	right    vec3.T
	up       vec3.T
	position vec3.T
	rotation vec3.T
	scale    vec3.T
}

// NewTransform creates a new 3D transform
func New(position, rotation, scale vec3.T) *T {
	t := &T{
		World:    mat4.Ident(),
		Local:    mat4.Ident(),
		position: position,
		rotation: rotation,
		scale:    scale,
		right:    vec3.UnitX,
		up:       vec3.UnitY,
		forward:  vec3.UnitZN,
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
	m := mat4.Transform(t.position, t.rotation, t.scale)

	// Update parent -> local transformation matrix
	t.Local = m

	// Update local -> world matrix
	if parent != nil {
		t.World = m.Mul(parent)
	} else {
		t.World = m
	}

	// Grab axis vectors from transformation matrix
	t.up = t.World.Up()
	t.right = t.World.Right()
	t.forward = t.World.Forward()
}

// Translate this transform by the given offset
func (t *T) Translate(offset vec3.T) {
	t.position = t.position.Add(offset)
}

// TransformPoint transforms a world point into this coordinate system
func (t *T) TransformPoint(point vec3.T) vec3.T {
	return t.World.TransformPoint(point)
}

// TransformDir transforms a world direction vector into this coordinate system
func (t *T) TransformDir(dir vec3.T) vec3.T {
	return t.World.TransformDir(dir)
}

func (t *T) Forward() vec3.T      { return t.forward }
func (t *T) Right() vec3.T        { return t.right }
func (t *T) Up() vec3.T           { return t.up }
func (t *T) Position() vec3.T     { return t.position }
func (t *T) Rotation() vec3.T     { return t.rotation }
func (t *T) Scale() vec3.T        { return t.scale }
func (t *T) SetPosition(p vec3.T) { t.position = p }
func (t *T) SetRotation(r vec3.T) { t.rotation = r }
func (t *T) SetScale(s vec3.T)    { t.scale = s }
