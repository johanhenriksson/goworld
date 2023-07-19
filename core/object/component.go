package object

import (
	"reflect"

	"github.com/johanhenriksson/goworld/core/transform"
)

type Component interface {
	// ID returns a unique identifier for this object.
	ID() uint

	// Name is used to identify the object within the scene.
	Name() string

	// Parent returns the parent of this object, or nil
	Parent() Object

	// Transform returns the object transform
	Transform() transform.T

	// Active indicates whether the object is active in the scene or not.
	// E.g. the object/component and all its parents are enabled and active.
	Active() bool

	// Enabled indicates whether the object is currently enabled or not.
	// Note that the object can still be inactive if an ancestor is disabled.
	Enabled() bool

	// Update the object. Called on every frame.
	Update(Component, float32)

	// Destroy the object
	Destroy()

	setName(string)
	setParent(Object)
	setEnabled(bool) bool
	setActive(bool) bool
}

type base struct {
	id      uint
	name    string
	enabled bool
	active  bool
	parent  Object
}

func emptyBase(name string) *base {
	return &base{
		id:      ID(),
		name:    name,
		enabled: true,
		active:  false,
	}
}

var componentType = reflect.TypeOf((*Component)(nil)).Elem()

func NewComponent[K Component](cmp K) K {
	t := reflect.TypeOf(cmp).Elem()
	v := reflect.ValueOf(cmp).Elem()

	// find & initialize base component
	baseIdx := -1
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.Anonymous {
			// only anonymous fields are considered
			continue
		}
		if !field.IsExported() {
			// only exported fields can be base fields
			continue
		}

		value := v.Field(i)
		if field.Type == componentType {
			// the components directly extends the base component
			// if its nil, create a new empty component base
			if value.IsZero() {
				base := emptyBase(t.Name())
				value.Set(reflect.ValueOf(base))
			}
		} else if _, isComponent := value.Interface().(Component); isComponent {
			// this object extends some other non-base object
		} else {
			// its not an object, move on
			continue
		}

		baseIdx = i
	}
	if baseIdx < 0 {
		panic("struct does not embed a Component")
	}

	return cmp
}

func (b *base) ID() uint {
	return b.id
}

func (b *base) Update(scene Component, dt float32) {
}

func (b *base) Transform() transform.T {
	if b.parent == nil {
		return transform.Identity()
	}
	return b.parent.Transform()
}

func (b *base) Active() bool { return b.active }
func (b *base) setActive(active bool) bool {
	prev := b.active
	b.active = active
	return prev
}

func (b *base) Enabled() bool { return b.enabled }
func (b *base) setEnabled(enabled bool) bool {
	prev := b.enabled
	b.enabled = enabled
	return prev
}

func (b *base) Parent() Object { return b.parent }
func (b *base) setParent(p Object) {
	if b.parent == p {
		return
	}
	b.parent = p
}

func (b *base) setName(n string) { b.name = n }
func (b *base) Name() string     { return b.name }
func (b *base) String() string   { return b.Name() }

func (o *base) Destroy() {
	if o.parent != nil {
		o.parent.detach(o)
	}
}
