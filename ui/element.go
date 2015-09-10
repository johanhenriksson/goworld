package ui;

import (
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"
)

type Element struct {
    Width       float64
    Height      float64
    Parent      *Element
    Children    []*Element
    Transform   *engine.Transform
    Material    *render.Material
}

func (e *Element) AddChild(child *Element) {
    e.Children = append(e.Children, child)
    child.Parent = e
}

func (e *Element) Draw() {
}
