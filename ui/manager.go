package ui;

import (
    "github.com/johanhenriksson/goworld/engine"
)

/** Main UI manager. Handles routing of events and drawing the UI. */
type Manager struct {
    Window      *engine.Window

    /** A top level element that represents the screen area. */
    Screen      Element

    /** Projection matrix - orthographic */
    /* Shaders */
}

func NewManager(wnd *engine.Window) *Manager {
    m := &Manager {
        Window: wnd,
    }
    return m
}

func (m *Manager) Draw() {
}

func (m *Manager) NewElement() *Element {
    e := &Element {
        Width: 0,
        Height: 0,
        Transform: engine.CreateTransform(0,0,0),
        Children: []*Element{},
    }
    return e
}

func (m *Manager) NewRect() *Rect {
    r := &Rect {
        Element: m.NewElement(),
    }
    return r
}
