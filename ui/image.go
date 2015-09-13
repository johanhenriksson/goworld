package ui

import (
    "github.com/johanhenriksson/goworld/render"
)

type Image struct {
    *Element
    Image   *render.Texture
    quad    *QuadImg
}

func (m *Manager) NewImage(image *render.Texture, x, y, w, h, z float32) *Image {
    el := m.NewElement(x,y,w,h,z)
    mat := render.LoadMaterial("assets/materials/ui_texture.json")
    mat.AddTexture(0, image)
    img := &Image {
        Element: el,
        Image: image,
        quad: NewQuadImg(mat, w, h, z),
    }
    return img
}

func (r *Image) Draw(args DrawArgs) {
    args.Transform = r.Element.Transform.Matrix.Mul4(args.Transform) //args.Transform.Mul4(r.Element.Transform.Matrix)
    r.quad.Draw(args)
    for _, el := range r.Element.children {
        el.Draw(args)
    }
}
