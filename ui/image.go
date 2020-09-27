package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
)

type Image struct {
	*Element
	Transparent bool
	Texture     *render.Texture
	Quad        *geometry.ImageQuad
}

func NewImage(texture *render.Texture, size vec2.T, invert bool, style Style) *Image {
	el := NewElement("Image", vec2.Zero, size, style)
	mat := assets.GetMaterial("ui_texture")
	mat.AddTexture("image", texture)
	return &Image{
		Element:     el,
		Texture:     texture,
		Quad:        geometry.NewImageQuad(mat, size, invert),
		Transparent: false,
	}
}

func NewDepthImage(texture *render.Texture, size vec2.T, invert bool) *Image {
	el := NewElement("DepthImage", vec2.Zero, size, NoStyle)
	mat := assets.GetMaterial("ui_depth_texture")
	mat.AddTexture("image", texture)
	return &Image{
		Element: el,
		Texture: texture,
		Quad:    geometry.NewImageQuad(mat, size, invert),
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

func (r *Image) Resize(size vec2.T) vec2.T {
	if size.X != r.Width() || size.Y != r.Height() {
		r.Element.Resize(size)
		r.Quad.SetSize(size)
		u := math.Min(size.X/float32(r.Texture.Width), 1)
		v := math.Min(size.Y/float32(r.Texture.Height), 1)
		r.Quad.SetUV(u, v)
	}
	return r.Size
}

func (r *Image) Flow(available vec2.T) vec2.T {
	return r.Size
}
