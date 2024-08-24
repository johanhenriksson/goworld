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

// objectType caches a reference to Object's reflect.Type
var objectType = reflect.TypeOf((*Object)(nil)).Elem()     // Object
var baseObjectType = reflect.TypeOf((*object)(nil)).Elem() // object

type object struct {
	component
	transform transform.T
	children  []Component
}

func emptyObject(pool Pool, name string) *object {
	return &object{
		transform: transform.Identity(),
		component: emptyComponent(name),
	}
}

// Empty creates a new, empty object.
func Empty(pool Pool, name string) Object {
	obj := emptyObject(pool, name)
	pool.assign(obj)
	return obj
}

func NewObject[K Object](pool Pool, name string, obj K) K {
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
				base := Empty(pool, name)
				value.Set(reflect.ValueOf(base))
			}
		} else if _, isObject := value.Interface().(Object); isObject {
			// this object extends some other non-base object
			if value.IsZero() {
				panic("base object is not initialized")
			}
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

	pool.assign(obj)

	// add Component fields as children
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if i == baseIdx {
			continue
		}

		if field.Anonymous {
			panic("multiple embeds are not allowed")
		}

		if !field.IsExported() {
			continue
		}

		// all uninitialized fields are ignored since they cant contain valid component references
		value := v.Field(i)
		if value.IsZero() {
			// maybe this should not be allowed?
			continue
		}

		// the field is a component - add it to the object's children
		if child, ok := value.Interface().(Component); ok {
			if child.Parent() == nil {
				Attach(obj, child)
			} else {
				panic("embedded object/component is attached to another parent")
			}
		}
	}

	return obj
}

func (g *object) Transform() transform.T {
	return g.transform
}

func (o *object) setParent(parent Object) {
	// check for cycles
	ancestor := parent
	for ancestor != nil {
		if ancestor.ID() == o.ID() {
			panic("cyclical object hierarchies are not allowed")
		}
		ancestor = ancestor.Parent()
	}

	o.component.setParent(parent)
	if parent != nil {
		o.transform.SetParent(parent.Transform())
	} else {
		o.transform.SetParent(nil)
	}
}

func (g *object) Update(scene Component, dt float32) {
	// maybe hierarchical update is dumb?
	for _, child := range g.children {
		if child.Enabled() {
			child.Update(scene, dt)
		}
	}
}

func (g *object) Children() []Component {
	return g.children
}

func (g *object) attach(children ...Component) {
	for _, child := range children {
		g.attachIfNotChild(child)
	}
}

func (g *object) attachIfNotChild(child Component) {
	for _, existing := range g.children {
		if existing.ID() == child.ID() {
			return
		}
	}
	g.children = append(g.children, child)
}

func (g *object) detach(child Component) {
	for i, existing := range g.children {
		if existing.ID() == child.ID() {
			g.children = append(g.children[:i], g.children[i+1:]...)
			return
		}
	}
}

func (g *object) setActive(active bool) bool {
	wasActive := g.component.setActive(active)
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

func (g *object) KeyEvent(e keys.Event) {
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

func (g *object) MouseEvent(e mouse.Event) {
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

func (o *object) Destroy() {
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

	o.component.Destroy()
}
