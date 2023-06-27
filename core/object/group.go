package object

import (
	"log"
	"reflect"

	"github.com/johanhenriksson/goworld/core/input"
	"github.com/johanhenriksson/goworld/core/input/keys"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/transform"
)

type G interface {
	T
	input.Handler

	// Children returns a slice containing the objects children.
	Children() []T

	attach(...T)
	detach(T)
}

type group struct {
	base
	transform transform.T
	children  []T
}

// Empty creates a new, empty object.
func Empty(name string) G {
	return &group{
		base: base{
			id:        ID(),
			name:      name,
			enabled:   true,
			transform: transform.Identity(),
		},
	}
}

func Group[K G](name string, obj K) K {
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()

	// initialize group base
	init := false
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "G" {
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

func (g *group) Update(scene T, dt float32) {
	for _, child := range g.children {
		if child.Active() {
			child.Update(scene, dt)
		}
	}
}

func (g *group) Children() []T {
	return g.children
}

func (g *group) attach(children ...T) {
	for _, child := range children {
		g.attachIfNotChild(child)
	}
}

func (g *group) attachIfNotChild(child T) {
	for _, existing := range g.children {
		if existing.ID() == child.ID() {
			return
		}
	}
	g.children = append(g.children, child)
}

func (g *group) detach(child T) {
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
