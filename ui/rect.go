package ui

import (
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/geometry"
)

type Rect struct {
    *Element
    Color   render.Color
    quad    *geometry.Quad
}

func (m *Manager) NewRect(color render.Color, x, y, w, h, z float32) *Rect {
    el := m.NewElement(x,y,w,h,z)
    mat := render.LoadMaterial("assets/materials/ui_color.json")
    r := &Rect {
        Element: el,
        Color: color,
        quad: geometry.NewQuad(mat, color, w, h, z),
    }
    r.quad.SetBorderWidth(5)
    return r
}

func (r *Rect) Draw(args render.DrawArgs) {
    args.Transform = r.Element.Transform.Matrix.Mul4(args.Transform) //args.Transform.Mul4(r.Element.Transform.Matrix)
    r.quad.Draw(args)
    for _, el := range r.Element.children {
        el.Draw(args)
    }
}
