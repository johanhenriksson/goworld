package ui;

import (
    "github.com/johanhenriksson/goworld/engine"
)

type Element struct {
    Width       float64
    Height      float64
    Transform   *engine.Transform
    Parent      *Element
    Children    []*Element
}

func (e *Element) AddChild(child *Element) {
    e.Children = append(e.Children, child)
    child.Parent = e
}

func (e *Element) Draw() {
}
