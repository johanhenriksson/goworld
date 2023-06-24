package object

import (
	"log"
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
	Update(T, float32)

	// Destroy the object
	Destroy()

	ID() uint
	setName(string)
	setParent(T)
	attach(...T)
	detach(T)
}

type base struct {
	id        uint
	transform transform.T
	name      string
	enabled   bool
	parent    T
	children  []T
}

// Empty creates a new, empty object.
func Empty(name string) T {
	return &base{
		id:        ID(),
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
			if reflect.ValueOf(child) == reflect.Zero(field.Type) {
				log.Println(t.Name(), " child ", field.Name, " is nil")
				continue
			}
			// initialize recursively?
			if child.Parent() == nil {
				Attach(obj, child)
			}
		}
	}
	return obj
}

func (b *base) ID() uint {
	return b.id
}

func (b *base) Update(scene T, dt float32) {
	for _, child := range b.children {
		if child.Active() {
			child.Update(scene, dt)
		}
	}
}

func (b *base) Transform() transform.T {
	// todo: rewrite/refactor
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
	return b.children
}

func (b *base) attach(children ...T) {
	for _, child := range children {
		b.attachIfNotChild(child)
	}
}

func (b *base) attachIfNotChild(child T) {
	for _, existing := range b.children {
		if existing.ID() == child.ID() {
			return
		}
	}
	b.children = append(b.children, child)
}

func (b *base) detach(child T) {
	for i, existing := range b.children {
		if existing.ID() == child.ID() {
			b.children = append(b.children[:i], b.children[i+1:]...)
			return
		}
	}
}

func (b *base) setName(n string) { b.name = n }
func (b *base) Name() string     { return b.name }
func (b *base) String() string   { return b.Name() }

func (o *base) Destroy() {
	// iterate over a copy of the child slice, since it will be mutated
	// when the child detaches itself during destruction
	for _, child := range o.Children() {
		child.Destroy()
	}
	if o.parent != nil {
		o.parent.detach(o)
	}
}

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
