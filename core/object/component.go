package object

import (
	"reflect"

	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/transform"
)

type component struct {
	id     uint
	name   string
	parent T
}

var _ T = &component{}

func Component[K T](cmp K) K {
	t := reflect.TypeOf(cmp).Elem()
	v := reflect.ValueOf(cmp).Elem()

	// initialize object base
	init := false
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "T" {
			if v.Field(i).IsZero() {
				component := &component{
					id:   ID(),
					name: t.Name(),
				}
				v.Field(i).Set(reflect.ValueOf(component))
			}
			init = true
			break
		}
	}
	if !init {
		panic("struct does not appear to be a Component")
	}

	// ensure there are no child objects/components
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "T" {
			continue
		}
		if !field.IsExported() {
			continue
		}
		if _, ok := v.Field(i).Interface().(T); ok {
			panic("Components cant have children")
		}
	}
	return cmp
}

func (b *component) ID() uint {
	return b.id
}

func (b *component) Active() bool     { return true }
func (b *component) SetActive(v bool) {}
func (b *component) Children() []T    { return nil }
func (b *component) Destroy()         {}

func (b *component) Update(scene T, dt float32) {
}

func (b *component) Transform() transform.T {
	return b.parent.Transform()
}

func (b *component) setName(n string) { b.name = n }
func (b *component) Name() string     { return b.name }
func (b *component) String() string   { return b.Name() }

func (b *component) Parent() T { return b.parent }
func (b *component) setParent(p T) {
	b.parent = p
}

func (b *component) attach(children ...T) {
	panic("cant attach children to components")
}

func (b *component) detach(child T) {
	panic("cant detach children to components")
}

func (o *component) KeyEvent(e keys.Event) {}

func (o *component) MouseEvent(e mouse.Event) {}
