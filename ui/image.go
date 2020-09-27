package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/render"
)

type Image struct {
	*Element
	Transparent bool
	Texture     *render.Texture
	Quad        *geometry.ImageQuad
}

func NewImage(texture *render.Texture, w, h float32, invert bool, style Style) *Image {
	el := NewElement("Image", 0, 0, w, h, style)
	mat := assets.GetMaterial("ui_texture")
	mat.AddTexture("image", texture)
	return &Image{
		Element:     el,
		Texture:     texture,
		Quad:        geometry.NewImageQuad(mat, w, h, invert),
		Transparent: false,
	}
}

func NewDepthImage(texture *render.Texture, w, h float32, invert bool) *Image {
	el := NewElement("DepthImage", 0, 0, w, h, NoStyle)
	mat := assets.GetMaterial("ui_depth_texture")
	mat.AddTexture("image", texture)
	return &Image{
		Element: el,
		Texture: texture,
		Quad:    geometry.NewImageQuad(mat, w, h, invert),
	}
}

func (r *Image) Draw(args render.DrawArgs) {
	args.Transform = r.Element.Transform.Matrix.Mul(&args.Transform) //args.Transform.Mul4(r.Element.Transform.Matrix)

	if r.Transparent {
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	} else {
		gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	}
	r.Quad.Material.Use()
	r.Quad.Material.RGBA("tint", r.Style.Color("color", render.White))
	r.Quad.Material.SetTexture("image", r.Texture)
	r.Quad.Draw(args)

	for _, el := range r.Element.children {
		el.Draw(args)
	}
}

func (r *Image) Resize(size Size) Size {
	if size.Width != r.Width() || size.Height != r.Height() {
		r.Element.Resize(size)
		r.Quad.SetSize(size.Width, size.Height)
		u := math.Min(size.Width/float32(r.Texture.Width), 1)
		v := math.Min(size.Height/float32(r.Texture.Height), 1)
		r.Quad.SetUV(u, v)
	}
	return r.Size
}

func (r *Image) Flow(available Size) Size {
	return r.Size
}
