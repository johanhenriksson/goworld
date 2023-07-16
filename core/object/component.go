package object

import (
	"log"
	"reflect"

	"github.com/johanhenriksson/goworld/core/transform"
)

type Component interface {
	// Name is used to identify the object within the scene.
	Name() string

	// Parent returns the parent of this object, or nil
	Parent() G

	// Transform returns the object transform
	Transform() transform.T

	// Active indicates whether the object is currently enabled or not.
	Active() bool

	// SetActive enables or disables the object
	SetActive(bool)

	// Update the object. Called on every frame.
	Update(Component, float32)

	// Destroy the object
	Destroy()

	ID() uint
	setName(string)
	setParent(G)
}

type base struct {
	id      uint
	name    string
	enabled bool
	parent  G
}

func emptyBase(name string) *base {
	return &base{
		id:      ID(),
		name:    name,
		enabled: true,
	}
}

func NewComponent[K Component](obj K) K {
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()

	// initialize object base
	init := false
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "Component" {
			if v.Field(i).IsZero() {
				base := emptyBase(t.Name())
				v.Field(i).Set(reflect.ValueOf(base))
			}
			init = true
			break
		}
	}
	if !init {
		// todo: does this even matter?
		// this forces extending structs to be named Component as well
		panic("struct does not appear to be a Component")
	}

	// add Object fields as children
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "Component" {
			continue
		}
		if !field.IsExported() {
			continue
		}
		if child, ok := v.Field(i).Interface().(Component); ok {
			if reflect.ValueOf(child) == reflect.Zero(field.Type) {
				log.Println(t.Name(), " child ", field.Name, " is nil")
				continue
			}
			// initialize recursively?
			panic("only groups can have children")
		}
	}
	return obj
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

func (b *base) Active() bool { return b.enabled }

func (b *base) SetActive(active bool) {
	if b.enabled && !active {
		// disable
		// if attached, raise OnDeactivate()
		b.enabled = false
	}
	if !b.enabled && active {
		// enable
		// if attached, raise OnActivate()
		b.enabled = true
	}
}

func (b *base) Parent() G { return b.parent }
func (b *base) setParent(p G) {
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
