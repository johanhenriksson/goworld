package physics

import (
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Transform represents a 3D transformation for a physics object.
// Implements the transform.T interface, but does not consider parent transformations.
type Transform struct {
	position vec3.T
	scale    vec3.T
	rotation quat.T

	matrix  mat4.T
	right   vec3.T
	up      vec3.T
	forward vec3.T

	inv   *mat4.T
	dirty bool
}

var _ transform.T = &Transform{}

func newTransform(position vec3.T, rotation quat.T, scale vec3.T) *Transform {
	t := &Transform{
		matrix:   mat4.Ident(),
		position: position,
		rotation: rotation,
		scale:    scale,
	}
	t.Recalculate(nil)
	return t
}

// Identity returns a new transform that does nothing.
func identity() *Transform {
	return newTransform(vec3.Zero, quat.Ident(), vec3.One)
}

// Update transform matrix and its right/up/forward vectors
func (t *Transform) Recalculate(parent transform.T) {
	if !t.dirty {
		// no parent, no change -> nothing to do
		// this will be common for floating objects, i.e. rigidbodies
		return
	}

	// calculate basis vectors
	t.right = t.rotation.Rotate(vec3.Right)
	t.up = t.rotation.Rotate(vec3.Up)
	t.forward = t.rotation.Rotate(vec3.Forward)

	// apply scaling
	x := t.right.Scaled(t.scale.X)
	y := t.up.Scaled(t.scale.Y)
	z := t.forward.Scaled(t.scale.Z)

	// create transformation matrix
	p := t.position
	t.matrix = mat4.T{
		x.X, x.Y, x.Z, 0,
		y.X, y.Y, y.Z, 0,
		z.X, z.Y, z.Z, 0,
		p.X, p.Y, p.Z, 1,
	}

	// mark as clean
	t.dirty = false

	// clear inversion cache
	t.inv = nil
}

func (t *Transform) inverse() *mat4.T {
	if t.inv == nil {
		inv := t.matrix.Invert()
		t.inv = &inv
	}
	return t.inv
}

// Translate this transform by the given offset
func (t *Transform) Translate(offset vec3.T) {
	t.position = t.position.Add(offset)
}

func (t *Transform) Project(point vec3.T) vec3.T {
	return t.matrix.TransformPoint(point)
}

func (t *Transform) Unproject(point vec3.T) vec3.T {
	return t.inverse().TransformPoint(point)
}

func (t *Transform) ProjectDir(dir vec3.T) vec3.T {
	return t.matrix.TransformDir(dir)
}

func (t *Transform) UnprojectDir(dir vec3.T) vec3.T {
	return t.inverse().TransformDir(dir)
}

func (t *Transform) WorldPosition() vec3.T       { return t.position }
func (t *Transform) SetWorldPosition(pos vec3.T) { t.position = pos; t.dirty = true }

func (t *Transform) WorldRotation() quat.T       { return t.rotation }
func (t *Transform) SetWorldRotation(rot quat.T) { t.rotation = rot; t.dirty = true }

func (t *Transform) Matrix() mat4.T  { return t.matrix }
func (t *Transform) Right() vec3.T   { return t.right }
func (t *Transform) Up() vec3.T      { return t.up }
func (t *Transform) Forward() vec3.T { return t.forward }

func (t *Transform) Position() vec3.T     { return t.position }
func (t *Transform) Rotation() quat.T     { return t.rotation }
func (t *Transform) Scale() vec3.T        { return t.scale }
func (t *Transform) SetPosition(p vec3.T) { t.position = p; t.dirty = true }
func (t *Transform) SetRotation(r quat.T) { t.rotation = r; t.dirty = true }
func (t *Transform) SetScale(s vec3.T)    { t.scale = s; t.dirty = true }
