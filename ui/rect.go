package ui

import (
    "github.com/johanhenriksson/goworld/render"
)

type Rect struct {
    *Element
}

func (m *Manager) NewRect(x, y, w, h float32) *Rect {
    el := m.NewElement(x,y,w,h)
    r := &Rect {
        Element: el,
    }
    el.Material = render.LoadMaterial("assets/materials/ui_color.json")
    /* TODO: set viewport matrix on material */
    return r
}

func (r *Rect) Draw(args DrawArgs) {
    sh := r.Element.Material.Shader
    r.Element.Material.Use()

    sh.Matrix4f("viewport", &args.Viewport[0])
    sh.Matrix4f("model", &r.Transform.Matrix[0])

    r.Element.Draw(args)
}
