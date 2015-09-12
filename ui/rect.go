package ui

import (
    "github.com/johanhenriksson/goworld/render"
)

type Rect struct {
    *Element
    Color   Color
    quad    *Quad
}

func (m *Manager) NewRect(color Color, x, y, w, h float32) *Rect {
    el := m.NewElement(x,y,w,h)
    mat := render.LoadMaterial("assets/materials/ui_color.json")
    r := &Rect {
        Element: el,
        Color: color,
        quad: NewQuad(mat, color, w, h, 0, 0,1,0,1),
    }
    return r
}

func (r *Rect) Draw(args DrawArgs) {
    args.Transform = args.Transform.Mul4(r.Element.Transform.Matrix)
    r.quad.Draw(args)
    r.Element.Draw(args)
}
