package transform

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type T interface {
	Recalculate(T)

	Forward() vec3.T
	Right() vec3.T
	Up() vec3.T

	Position() vec3.T
	SetPosition(vec3.T)

	Rotation() quat.T
	SetRotation(quat.T)

	Scale() vec3.T
	SetScale(vec3.T)

	Matrix() mat4.T

	ProjectDir(dir vec3.T) vec3.T

	// Project local coordinates to world coordinates
	Project(point vec3.T) vec3.T

	// Unproject world coordinates to local coordinates
	Unproject(point vec3.T) vec3.T

	UnprojectDir(dir vec3.T) vec3.T

	WorldPosition() vec3.T
	SetWorldPosition(vec3.T)
}

// Transform represents a 3D transformation
type transform struct {
	position vec3.T
	scale    vec3.T
	rotation quat.T

	matrix  mat4.T
	right   vec3.T
	up      vec3.T
	forward vec3.T
}

// NewTransform creates a new 3D transform
func New(position vec3.T, rotation quat.T, scale vec3.T) T {
	t := &transform{
		matrix:   mat4.Ident(),
		position: position,
		rotation: rotation,
		scale:    scale,
	}
	t.Recalculate(nil)
	return t
}

// Identity returns a new transform that does nothing.
func Identity() T {
	return New(vec3.Zero, quat.Ident(), vec3.One)
}

// Update transform matrix and its right/up/forward vectors
func (t *transform) Recalculate(parent T) {
	position := t.position
	rotation := t.rotation
	scale := t.scale

	if parent != nil {
		scale = scale.Mul(parent.Scale())

		rotation = rotation.Mul(parent.Rotation())

		position = parent.Rotation().Rotate(parent.Scale().Mul(position))
		position = parent.Position().Add(position)
	}

	// calculate basis vectors
	t.right = rotation.Rotate(vec3.Right)
	t.up = rotation.Rotate(vec3.Up)
	t.forward = rotation.Rotate(vec3.Forward)

	// apply scaling
	x := t.right.Scaled(scale.X)
	y := t.up.Scaled(scale.Y)
	z := t.forward.Scaled(scale.Z)

	// create transformation matrix
	p := position
	t.matrix = mat4.T{
		x.X, x.Y, x.Z, 0,
		y.X, y.Y, y.Z, 0,
		z.X, z.Y, z.Z, 0,
		p.X, p.Y, p.Z, 1,
	}
}

// Translate this transform by the given offset
func (t *transform) Translate(offset vec3.T) {
	t.position = t.position.Add(offset)
}

func (t *transform) Project(point vec3.T) vec3.T {
	return t.matrix.TransformPoint(point)
}

func (t *transform) Unproject(point vec3.T) vec3.T {
	inv := t.matrix.Invert()
	return inv.TransformPoint(point)
}

func (t *transform) ProjectDir(dir vec3.T) vec3.T {
	return t.matrix.TransformDir(dir)
}

func (t *transform) UnprojectDir(dir vec3.T) vec3.T {
	inv := t.matrix.Invert()
	return inv.TransformDir(dir)
}

func (t *transform) WorldPosition() vec3.T {
	return t.matrix.Origin()
}

func (t *transform) SetWorldPosition(wp vec3.T) {
	offset := t.Unproject(wp)
	t.SetPosition(t.position.Add(offset))
}

func (t *transform) Matrix() mat4.T  { return t.matrix }
func (t *transform) Right() vec3.T   { return t.right }
func (t *transform) Up() vec3.T      { return t.up }
func (t *transform) Forward() vec3.T { return t.forward }

func (t *transform) Position() vec3.T     { return t.position }
func (t *transform) Rotation() quat.T     { return t.rotation }
func (t *transform) Scale() vec3.T        { return t.scale }
func (t *transform) SetPosition(p vec3.T) { t.position = p }
func (t *transform) SetRotation(r quat.T) { t.rotation = r }
func (t *transform) SetScale(s vec3.T)    { t.scale = s }

func Matrix(position vec3.T, rotation quat.T, scale vec3.T) mat4.T {
	x := rotation.Rotate(vec3.Right)
	y := rotation.Rotate(vec3.Up)
	z := rotation.Rotate(vec3.Forward)

	x.Scale(scale.X)
	y.Scale(scale.Y)
	z.Scale(scale.Z)

	p := position
	return mat4.T{
		x.X, x.Y, x.Z, 0,
		y.X, y.Y, y.Z, 0,
		z.X, z.Y, z.Z, 0,
		p.X, p.Y, p.Z, 1,
	}
}
