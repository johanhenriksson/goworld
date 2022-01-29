package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

type Image struct {
	*Element
	Transparent bool
	Texture     texture.T
	mesh        *geometry.Rect
}

func NewImage(texture texture.T, size vec2.T, invert bool, style Style) *Image {
	el := NewElement("Image", vec2.Zero, size, style)
	mat := assets.GetMaterial("ui_texture")
	mat.Texture("image", texture)
	rect := geometry.NewRect(mat, size)
	rect.Invert = invert
	return &Image{
		Element:     el,
		Texture:     texture,
		Transparent: false,
		mesh:        rect,
	}
}

func NewDepthImage(texture texture.T, size vec2.T, invert bool) *Image {
	el := NewElement("DepthImage", vec2.Zero, size, NoStyle)
	mat := assets.GetMaterial("ui_texture")
	mat.Texture("image", texture)
	rect := geometry.NewRect(mat, size)
	rect.Invert = invert
	rect.Depth = true
	return &Image{
		Element: el,
		Texture: texture,
		mesh:    rect,
	}
}

func (r *Image) Draw(args render.Args) {
	args.Transform = r.Element.Transform.Matrix.Mul(&args.Transform) //args.Transform.Mul4(r.Element.Transform.Matrix)

	if r.Transparent {
		render.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	} else {
		render.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	}
	r.mesh.Material.Use()
	r.mesh.Material.RGBA("tint", r.Style.Color("color", color.White))
	r.mesh.Material.Texture("image", r.Texture)
	r.mesh.Material.Bool("depth", r.mesh.Depth)
	r.mesh.Material.Bool("invert", r.mesh.Invert)
	r.mesh.Draw(args)

	for _, el := range r.Element.children {
		el.Draw(args)
	}
}

func (r *Image) Resize(size vec2.T) vec2.T {
	if size.X != r.Width() || size.Y != r.Height() {
		r.Element.Resize(size)
		r.mesh.SetSize(size)
	}
	return r.Size
}

func (r *Image) Flow(available vec2.T) vec2.T {
	return r.Size
}
