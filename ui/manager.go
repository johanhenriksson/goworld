package ui;

import (
    "github.com/johanhenriksson/goworld/render"
    mgl "github.com/go-gl/mathgl/mgl32"
)

/** Main UI manager. Handles routing of events and drawing the UI. */
type Manager struct {
    /** Projection matrix - orthographic */
    Viewport    mgl.Mat4

    Children    []render.Drawable
}

func NewManager(width, height float32) *Manager {
    m := &Manager {
        Viewport: mgl.Ortho(0, width, 0, height, 1000, -1000),
        Children: []render.Drawable{},
    }
    return m
}

func (m *Manager) Append(child render.Drawable) {
    m.Children = append(m.Children, child)
}

func (m *Manager) Draw() {
    /* create draw event args */
    args := render.DrawArgs {
        Projection: m.Viewport,
        View: mgl.Ident4(), // unused
        Transform: mgl.Ident4(),
    }

    for _, el := range m.Children {
        el.Draw(args)
    }
}

