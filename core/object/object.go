package object

import (
	"log"
	"reflect"

	"github.com/johanhenriksson/goworld/core/input"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/transform"
)

type Object interface {
	Component
	input.Handler

	// Children returns a slice containing the objects children.
	Children() []Component

	attach(...Component)
	detach(Component)
}

type group struct {
	base
	transform transform.T
	children  []Component
}

// Empty creates a new, empty object.
func Empty(name string) Object {
	return &group{
		base: base{
			id:      ID(),
			name:    name,
			enabled: true,
		},
		transform: transform.Identity(),
	}
}

func New[K Object](name string, obj K) K {
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()

	// initialize group base
	init := false
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "Object" {
			if v.Field(i).IsZero() {
				base := Empty(name)
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
		if field.Name == "Object" {
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
			if child.Parent() == nil {
				Attach(obj, child)
			}
		}
	}
	return obj
}

func (g *group) Transform() transform.T {
	// todo: rewrite/refactor
	var pt transform.T = nil
	if g.parent != nil {
		pt = g.parent.Transform()
	}
	g.transform.Recalculate(pt)

	return g.transform
}

func (g *group) Update(scene Component, dt float32) {
	for _, child := range g.children {
		if child.Active() {
			child.Update(scene, dt)
		}
	}
}

func (g *group) Children() []Component {
	return g.children
}

func (g *group) attach(children ...Component) {
	for _, child := range children {
		g.attachIfNotChild(child)
	}
}

func (g *group) attachIfNotChild(child Component) {
	for _, existing := range g.children {
		if existing.ID() == child.ID() {
			return
		}
	}
	g.children = append(g.children, child)
}

func (g *group) detach(child Component) {
	for i, existing := range g.children {
		if existing.ID() == child.ID() {
			g.children = append(g.children[:i], g.children[i+1:]...)
			return
		}
	}
}

func (g *group) KeyEvent(e keys.Event) {
	for _, child := range g.children {
		if !child.Active() {
			continue
		}
		if handler, ok := child.(input.KeyHandler); ok {
			handler.KeyEvent(e)
			if e.Handled() {
				return
			}
		}
	}
}

func (g *group) MouseEvent(e mouse.Event) {
	for _, child := range g.children {
		if !child.Active() {
			continue
		}
		if handler, ok := child.(input.MouseHandler); ok {
			handler.MouseEvent(e)
			if e.Handled() {
				return
			}
		}
	}
}

func (o *group) Destroy() {
	// iterate over a copy of the child slice, since it will be mutated
	// when the child detaches itself during destruction
	children := make([]Component, len(o.Children()))
	copy(children, o.Children()[:])

	for _, child := range o.Children() {
		child.Destroy()
	}

	if o.parent != nil {
		o.parent.detach(o)
	}
}
