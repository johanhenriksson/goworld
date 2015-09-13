package ui

import (
    "github.com/johanhenriksson/goworld/render"
)

type Rect struct {
    *Element
    Color   Color
    quad    *Quad
}

func (m *Manager) NewRect(color Color, x, y, w, h, z float32) *Rect {
    el := m.NewElement(x,y,w,h,z)
    mat := render.LoadMaterial("assets/materials/ui_color.json")
    r := &Rect {
        Element: el,
        Color: color,
        quad: NewQuad(mat, color, w, h, z, 0,1,0,1),
    }
    r.quad.SetBorderWidth(10)
    return r
}

func (r *Rect) Draw(args DrawArgs) {
    args.Transform = r.Element.Transform.Matrix.Mul4(args.Transform) //args.Transform.Mul4(r.Element.Transform.Matrix)
    r.quad.Draw(args)
    for _, el := range r.Element.children {
        el.Draw(args)
    }
}
