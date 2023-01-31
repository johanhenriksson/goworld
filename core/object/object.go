package object

import (
	"reflect"

	"github.com/johanhenriksson/goworld/core/input"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/transform"
)

type T interface {
	input.Handler

	// Name is used to identify the object within the scene.
	Name() string

	// Parent returns the parent of this object, or nil
	Parent() T

	// Children returns a slice containing the objects children.
	Children() []T

	// Transform returns the object transform
	Transform() transform.T

	// Active indicates whether the object is currently enabled or not.
	Active() bool

	// SetActive enables or disables the object
	SetActive(bool)

	// Update the object. Called on every frame.
	Update(float32)

	setName(string)
	setParent(T)
	attach(...T)
	detach(T)
}

type base struct {
	transform transform.T
	name      string
	enabled   bool
	parent    T
	children  []T
}

// Empty creates a new, empty object.
func Empty(name string) T {
	return &base{
		name:      name,
		enabled:   true,
		transform: transform.Identity(),
	}
}

func New[K T](obj K) K {
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()

	// initialize object base
	init := false
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "T" {
			if v.Field(i).IsZero() {
				base := Empty(t.Name())
				v.Field(i).Set(reflect.ValueOf(base))
			}
			init = true
			break
		}
	}
	if !init {
		panic("struct does not appear to be an Object")
	}

	// add Object fields as children
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "T" {
			continue
		}
		if !field.IsExported() {
			continue
		}
		if child, ok := v.Field(i).Interface().(T); ok {
			// initialize recursively?
			if child.Parent() == nil {
				Attach(obj, child)
			}
		}
	}
	return obj
}

func (b *base) Update(dt float32) {
	for _, child := range b.children {
		if child.Active() {
			child.Update(dt)
		}
	}
}

func (b *base) Transform() transform.T {
	var pt transform.T = nil
	if b.parent != nil {
		pt = b.parent.Transform()
	}
	b.transform.Recalculate(pt)

	return b.transform
}

func (b *base) Active() bool     { return b.enabled }
func (b *base) SetActive(a bool) { b.enabled = a }

func (b *base) Parent() T { return b.parent }
func (b *base) setParent(p T) {
	if b.parent == p {
		return
	}
	b.parent = p
}

func (b *base) Children() []T {
	children := make([]T, len(b.children))
	copy(children, b.children)
	return children
}

func (b *base) attach(children ...T) {
	for _, child := range children {
		b.attachIfNotChild(child)
	}
}

func (b *base) attachIfNotChild(child T) {
	for _, existing := range b.children {
		if existing == child {
			return
		}
	}
	b.children = append(b.children, child)
}

func (b *base) detach(child T) {
	for i, existing := range b.children {
		if existing == child {
			b.children = append(b.children[:i], b.children[i+1:]...)
			return
		}
	}
}

func (b *base) setName(n string) { b.name = n }
func (b *base) Name() string     { return b.name }
func (b *base) String() string   { return b.Name() }

func (o *base) KeyEvent(e keys.Event) {
	for _, child := range o.children {
		if !child.Active() {
			continue
		}
		child.KeyEvent(e)
		if e.Handled() {
			return
		}
	}
}

func (o *base) MouseEvent(e mouse.Event) {
	for _, child := range o.children {
		if !child.Active() {
			continue
		}
		child.MouseEvent(e)
		if e.Handled() {
			return
		}
	}
}
