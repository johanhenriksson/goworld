package ui;

import (
    "github.com/johanhenriksson/goworld/engine"
)

type element_f struct {
    Type    string
    Width   float32
    Height  float32
}

type Element struct {
    width       float32
    height      float32
    Parent      Drawable
    Children    []Drawable
    Transform   *engine.Transform
}

func (m *Manager) NewElement(x, y, w, h float32) *Element {
    e := &Element {
        width: w,
        height: h,
        Transform: engine.CreateTransform(x,y,0),
        Children: []Drawable{},
    }
    return e
}

func (e *Element) Append(child *Element) {
    e.Children = append(e.Children, child)
}

func (e *Element) Remove(child *Element) {
    /* TODO: Implement */
    child.Parent = nil
}

func (e *Element) Draw(args DrawArgs) {
    /* Multiply transform to args */
    args.Transform = args.Transform.Mul4(e.Transform.Matrix)
    for _, el := range e.Children {
        el.Draw(args)
    }
}

func (e *Element) Width() float32 { return e.width; }
func (e *Element) Height() float32 { return e.height; }
