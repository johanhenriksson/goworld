package engine

import (
	"github.com/johanhenriksson/goworld/render"
	"reflect"
)

// Component is the general interface for scene object components.
type Component interface {
	Update(float32)
	Draw(render.DrawArgs)
	Object() *Object
	Type() reflect.Type
}

// ComponentBase holds the core information about a component, such as its type and parent object.
type ComponentBase struct {
	object *Object
	ctype  reflect.Type
}

// NewComponent creates a new base component and attaches it to a game object.
func NewComponent(parent *Object, component Component) *ComponentBase {
	c := &ComponentBase{
		object: parent,
		ctype:  reflect.TypeOf(component),
	}
	parent.Attach(component)
	return c
}

// Object returns the base game object
func (c *ComponentBase) Object() *Object {
	return c.object
}

// Type returns the component type
func (c *ComponentBase) Type() reflect.Type {
	return c.ctype
}

// GetComponent returns the first component of a given type on the parent object.
func (c *ComponentBase) GetComponent(component Component) (Component, bool) {
	return c.Object().GetComponent(component)
}
