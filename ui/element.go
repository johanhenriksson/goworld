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

func (e *Element) ZIndex() float32 { return e.z }
func (e *Element) Parent() render.Drawable { return e.parent }
func (e *Element) SetParent(parent render.Drawable) {
    e.parent = parent
}
func (e *Element) Children() []render.Drawable { return e.children }

func (e *Element) Append(child render.Drawable) {
    e.children = append(e.children, child)
}

func (e *Element) Remove(child render.Drawable) {
    /* TODO: Implement */
    //child.Parent = nil
}

func (e *Element) Draw(args render.DrawArgs) {
    /* Multiply transform to args */
    args.Transform = e.Transform.Matrix.Mul4(args.Transform) //args.Transform.Mul4(e.Transform.Matrix)
    for _, el := range e.children {
        el.Draw(args)
    }
}

func (e *Element) Width() float32 { return e.width; }
func (e *Element) Height() float32 { return e.height; }
