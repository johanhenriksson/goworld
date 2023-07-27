package transform

import (
	"github.com/johanhenriksson/goworld/core/events"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type T interface {
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

	WorldScale() vec3.T
	SetWorldScale(vec3.T)

	WorldRotation() quat.T
	SetWorldRotation(quat.T)

	Parent() T
	SetParent(T)
	OnChange() *events.Event[T]
}

// Transform represents a 3D transformation
type transform struct {
	position vec3.T
	scale    vec3.T
	rotation quat.T

	wposition vec3.T
	wscale    vec3.T
	wrotation quat.T

	matrix  mat4.T
	right   vec3.T
	up      vec3.T
	forward vec3.T

	inv     *mat4.T
	dirty   bool
	parent  T
	changed events.Event[T]
	unsub   func()
}

// NewTransform creates a new 3D transform
func New(position vec3.T, rotation quat.T, scale vec3.T) T {
	t := &transform{
		matrix:   mat4.Ident(),
		position: position,
		rotation: rotation,
		scale:    scale,
		dirty:    true,
	}
	t.refresh()
	return t
}

// Identity returns a new transform that does nothing.
func Identity() T {
	return New(vec3.Zero, quat.Ident(), vec3.One)
}

func (t *transform) Parent() T { return t.parent }
func (t *transform) SetParent(parent T) {
	// check for cycles
	ancestor := parent
	for ancestor != nil {
		if ancestor.(T) == t {
			panic("cyclical transform hierarchies are not allowed")
		}
		ancestor = ancestor.Parent()
	}

	// todo: we might want to maintain world transform when attaching/detaching

	// detach from previous parent (if any)
	if t.unsub != nil {
		// unsub
		t.unsub()
		t.unsub = nil
	}

	t.parent = parent

	// attach to new parent (if any)
	if t.parent != nil {
		t.unsub = t.parent.OnChange().Subscribe(func(parent T) {
			// mark as dirty on parent change
			t.invalidate()
		})
	}

	t.invalidate()
}

func (t *transform) OnChange() *events.Event[T] {
	return &t.changed
}

func (t *transform) invalidate() {
	t.dirty = true
	t.changed.Emit(t)
}

// Update transform matrix and its right/up/forward vectors
func (t *transform) refresh() {
	if !t.dirty {
		return
	}

	position := t.position
	rotation := t.rotation
	scale := t.scale

	if t.parent != nil {
		scale = scale.Mul(t.parent.WorldScale())

		rotation = rotation.Mul(t.parent.WorldRotation())

		position = t.parent.WorldRotation().Rotate(t.parent.WorldScale().Mul(position))
		position = t.parent.WorldPosition().Add(position)
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

	// save world transforms
	t.wposition = position
	t.wscale = scale
	t.wrotation = rotation

	// mark as clean
	t.dirty = false

	// clear inversion cache
	t.inv = nil
}

func (t *transform) inverse() *mat4.T {
	t.refresh()
	if t.inv == nil {
		inv := t.matrix.Invert()
		t.inv = &inv
	}
	return t.inv
}

func (t *transform) Project(point vec3.T) vec3.T {
	t.refresh()
	return t.matrix.TransformPoint(point)
}

func (t *transform) Unproject(point vec3.T) vec3.T {
	return t.inverse().TransformPoint(point)
}

func (t *transform) ProjectDir(dir vec3.T) vec3.T {
	t.refresh()
	return t.matrix.TransformDir(dir)
}

func (t *transform) UnprojectDir(dir vec3.T) vec3.T {
	return t.inverse().TransformDir(dir)
}

func (t *transform) WorldPosition() vec3.T { t.refresh(); return t.wposition }
func (t *transform) WorldScale() vec3.T    { t.refresh(); return t.wscale }
func (t *transform) WorldRotation() quat.T { t.refresh(); return t.wrotation }

func (t *transform) SetWorldPosition(wp vec3.T) {
	t.refresh()
	if t.parent != nil {
		t.position = t.parent.Unproject(wp)
	} else {
		t.position = wp
	}
	t.invalidate()
}

func (t *transform) SetWorldScale(wscale vec3.T) {
	t.refresh()
	if t.parent != nil {
		// world_scale = parent_scale * local_scale
		// local_scale = world_scale / parent_scale
		t.scale = wscale.Div(t.parent.WorldScale())
	} else {
		t.scale = wscale
	}
	t.invalidate()
}

func (t *transform) SetWorldRotation(rot quat.T) {
	t.refresh()
	if t.parent != nil {
		t.rotation = rot.Mul(t.parent.WorldRotation().Inverse())
	} else {
		t.rotation = rot
	}
	t.invalidate()
}

func (t *transform) Matrix() mat4.T  { t.refresh(); return t.matrix }
func (t *transform) Right() vec3.T   { t.refresh(); return t.right }
func (t *transform) Up() vec3.T      { t.refresh(); return t.up }
func (t *transform) Forward() vec3.T { t.refresh(); return t.forward }

func (t *transform) Position() vec3.T     { t.refresh(); return t.position }
func (t *transform) Rotation() quat.T     { t.refresh(); return t.rotation }
func (t *transform) Scale() vec3.T        { t.refresh(); return t.scale }
func (t *transform) SetPosition(p vec3.T) { t.position = p; t.invalidate() }
func (t *transform) SetRotation(r quat.T) { t.rotation = r; t.invalidate() }
func (t *transform) SetScale(s vec3.T)    { t.scale = s; t.invalidate() }

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
