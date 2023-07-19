package object

import (
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

var objectType = reflect.TypeOf((*Object)(nil)).Elem()

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

	// find & initialize base object
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
		if field.Type == objectType {
			// the object directly extends the base object
			// if its nil, create a new empty object base
			if value.IsZero() {
				base := Empty(name)
				value.Set(reflect.ValueOf(base))
			}
		} else if _, isObject := value.Interface().(Object); isObject {
			// this object extends some other non-base object
		} else {
			// its not an object, move on
			continue
		}

		// if we already found a base field, the user has embedded multiple objects
		if baseIdx >= 0 {
			panic("struct embeds multiple Object types")
		}
		baseIdx = i
	}
	if baseIdx < 0 {
		panic("struct does not embed an Object")
	}

	// add Component fields as children
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if i == baseIdx {
			continue
		}

		if !field.IsExported() {
			continue
		}

		// all uninitialized fields are ignored since they cant contain valid component references
		value := v.Field(i)
		if value.IsZero() {
			continue
		}

		// the field contains a reference to an instantiated component
		// if its an orphan, add it to the object's children
		if child, ok := value.Interface().(Component); ok {
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
		if child.Enabled() {
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

func (g *group) setActive(active bool) bool {
	wasActive := g.base.setActive(active)
	if active {
		for _, child := range g.children {
			activate(child)
		}
	} else {
		for _, child := range g.children {
			deactivate(child)
		}
	}
	return wasActive
}

func (g *group) KeyEvent(e keys.Event) {
	for _, child := range g.children {
		if !child.Enabled() {
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
		if !child.Enabled() {
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
