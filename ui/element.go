package ui;

import (
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"
)

type Element struct {
    width       float32
    height      float32
    z           float32
    parent      render.Drawable
    children    []render.Drawable
    Transform   *engine.Transform
}

func (m *Manager) NewElement(x, y, w, h, z float32) *Element {
    e := &Element {
        width: w,
        height: h,
        children: []render.Drawable{},

        Transform: engine.CreateTransform(x,y,z),
    }
    return e
}

func (e *Element) ZIndex() float32 {
    // not sure how this is going to work yet
    // parents must be drawn underneath children (?)
    return e.z
}

/* Returns the parent element */
func (e *Element) Parent() render.Drawable {
    return e.parent
}

func (e *Element) SetParent(parent render.Drawable) {
    // TODO detach from current parent?
    e.parent = parent
}

func (e *Element) Children() []render.Drawable {
    return e.children
}

func (e *Element) Width() float32 {
    return e.width;
}

func (e *Element) Height() float32 {
    return e.height;
}

func (e *Element) Append(child render.Drawable) {
    e.children = append(e.children, child)
    // set parent?
}

func (e *Element) Remove(child render.Drawable) {
    // TODO Implement
    //child.Parent = nil
}

func (e *Element) Draw(args render.DrawArgs) {
    /* Multiply transform to args */
    args.Transform = e.Transform.Matrix.Mul4(args.Transform)
    for _, el := range e.children {
        el.Draw(args)
    }
}

