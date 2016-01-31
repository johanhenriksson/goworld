package ui

import (
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/geometry"
)

type Image struct {
    *Element
    Image   *render.Texture
    Quad    *geometry.ImageQuad
}

func (m *Manager) NewImage(image *render.Texture, x, y, w, h, z float32) *Image {
    el := m.NewElement(x,y,w,h,z)
    mat := render.LoadMaterial("assets/materials/ui_texture.json")
    mat.AddTexture("image", image)
    img := &Image {
        Element: el,
        Image: image,
        Quad: geometry.NewImageQuad(mat, w, h, z),
    }
    return img
}

func (r *Image) Draw(args render.DrawArgs) {
    args.Transform = r.Element.Transform.Matrix.Mul4(args.Transform) //args.Transform.Mul4(r.Element.Transform.Matrix)
    r.Quad.Draw(args)
    for _, el := range r.Element.children {
        el.Draw(args)
    }
}
