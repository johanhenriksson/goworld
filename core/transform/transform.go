package transform

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type T interface {
	Recalculate(T)

	Forward() vec3.T
	Right() vec3.T
	Up() vec3.T

	Position() vec3.T
	SetPosition(vec3.T)

	Rotation() vec3.T
	SetRotation(vec3.T)

	Scale() vec3.T
	SetScale(vec3.T)

	World() mat4.T
	Local() mat4.T

	ProjectDir(point vec3.T) vec3.T
	Project(point vec3.T) vec3.T
	Unproject(point vec3.T) vec3.T

	WorldPosition() vec3.T
}

// Transform represents a 3D transformation
type transform struct {
	world mat4.T
	local mat4.T

	forward  vec3.T
	right    vec3.T
	up       vec3.T
	position vec3.T
	rotation vec3.T
	scale    vec3.T
}

// NewTransform creates a new 3D transform
func New(position, rotation, scale vec3.T) T {
	t := &transform{
		world:    mat4.Ident(),
		local:    mat4.Ident(),
		position: position,
		rotation: rotation,
		scale:    scale,
		right:    vec3.UnitX,
		up:       vec3.UnitY,
		forward:  vec3.UnitZ,
	}
	t.Recalculate(nil)
	return t
}

// Identity returns a new transform that does nothing.
func Identity() T {
	return New(vec3.Zero, vec3.Zero, vec3.One)
}

// Update transform matrix and its right/up/forward vectors
func (t *transform) Recalculate(parent T) {
	// Update transform
	m := mat4.Transform(t.position, t.rotation, t.scale)

	// Update parent -> local transformation matrix
	t.local = m

	// Update local -> world matrix
	if parent != nil {
		pt := parent.World()
		t.world = pt.Mul(&t.local)
	} else {
		t.world = m
	}

	// Grab axis vectors from transformation matrix
	t.up = t.world.Up()
	t.right = t.world.Right()
	t.forward = t.world.Forward()
}

// Translate this transform by the given offset
func (t *transform) Translate(offset vec3.T) {
	t.position = t.position.Add(offset)
}

// TransformPoint transforms a world point into this coordinate system
func (t *transform) Project(point vec3.T) vec3.T {
	return t.world.TransformPoint(point)
}

func (t *transform) Unproject(point vec3.T) vec3.T {
	inv := t.world.Invert()
	return inv.TransformPoint(point)
}

// TransformDir transforms a world direction vector into this coordinate system
func (t *transform) ProjectDir(dir vec3.T) vec3.T {
	return t.world.TransformDir(dir)
}

func (t *transform) WorldPosition() vec3.T {
	return t.Project(vec3.Zero)
}

func (t *transform) World() mat4.T        { return t.world }
func (t *transform) Local() mat4.T        { return t.local }
func (t *transform) Forward() vec3.T      { return t.forward }
func (t *transform) Right() vec3.T        { return t.right }
func (t *transform) Up() vec3.T           { return t.up }
func (t *transform) Position() vec3.T     { return t.position }
func (t *transform) Rotation() vec3.T     { return t.rotation }
func (t *transform) Scale() vec3.T        { return t.scale }
func (t *transform) SetPosition(p vec3.T) { t.position = p }
func (t *transform) SetRotation(r vec3.T) { t.rotation = r }
func (t *transform) SetScale(s vec3.T)    { t.scale = s }
