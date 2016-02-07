package engine

import (
    "reflect"
    "github.com/johanhenriksson/goworld/render"
)

type Component interface {
    Update(float32)
    Draw(render.DrawArgs)
    Object() *Object
    Type() reflect.Type
}

type ComponentBase struct {
    object *Object
    ctype  reflect.Type
}

func NewComponent(parent *Object, component Component) *ComponentBase {
    c := &ComponentBase {
        object: parent,
        ctype: reflect.TypeOf(component),
    }
    parent.Attach(component)
    return c
}

func (c *ComponentBase) Object() *Object {
    return c.object
}

func (c *ComponentBase) Type() reflect.Type {
    return c.ctype
}

func (c *ComponentBase) GetComponent(component Component) (Component, bool) {
    return c.Object().GetComponent(component)
}