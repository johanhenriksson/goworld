package object

import (
	"reflect"

	"github.com/johanhenriksson/goworld/core/transform"
)

type Component interface {
	// ID returns a unique identifier for this object.
	ID() Handle

	// Pool returns the object pool this object belongs to.
	Pool() Pool

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

	// Update the component. Called on every frame.
	Update(Component, float32)

	setHandle(Pool, Handle)
	setName(string)
	setParent(Object)
	setEnabled(bool) bool
	setActive(bool) bool
	destroy()
}

type component struct {
	id      Handle
	name    string
	ctx     Pool
	enabled bool
	active  bool
	parent  Object
}

func emptyComponent(name string) component {
	return component{
		name:    name,
		enabled: true,
		active:  false,
	}
}

// componentType caches a reference to Component's reflect.Type
var componentType = reflect.TypeOf((*Component)(nil)).Elem()     // Component
var baseComponentType = reflect.TypeOf((*component)(nil)).Elem() // component

func NewComponent[K Component](pool Pool, cmp K) K {
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
				base := emptyComponent(t.Name())
				value.Set(reflect.ValueOf(&base))
			}
		} else if _, isComponent := value.Interface().(Component); isComponent {
			// this object extends some other non-base object
			if value.IsZero() {
				panic("base component is not initialized")
			}
		} else {
			// its not an object, move on
			continue
		}

		if baseIdx >= 0 {
			panic("struct embeds multiple Components")
		}
		baseIdx = i
	}
	if baseIdx < 0 {
		panic("struct does not embed a Component")
	}

	// assign context handle
	pool.assign(cmp)

	return cmp
}

func (b *component) ID() Handle { return b.id }
func (b *component) setHandle(pool Pool, id Handle) {
	b.id = id
	b.ctx = pool
}

func (b *component) Pool() Pool { return b.ctx }

func (b *component) Update(scene Component, dt float32) {}

func (b *component) Transform() transform.T {
	if b.parent == nil {
		return transform.Identity()
	}
	return b.parent.Transform()
}

func (b *component) Active() bool { return b.active }
func (b *component) setActive(active bool) bool {
	prev := b.active
	b.active = active
	return prev
}

func (b *component) Enabled() bool { return b.enabled }
func (b *component) setEnabled(enabled bool) bool {
	prev := b.enabled
	b.enabled = enabled
	return prev
}

func (b *component) Parent() Object     { return b.parent }
func (b *component) setParent(p Object) { b.parent = p }

func (b *component) setName(n string) { b.name = n }
func (b *component) Name() string     { return b.name }
func (b *component) String() string   { return b.Name() }

func (o *component) destroy() {
	if o.ctx != nil {
		o.ctx.release(o.id)
		o.ctx = nil
	}
}
