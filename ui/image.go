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
	Texture *render.Texture
	Quad    *geometry.ImageQuad
}

func NewImage(texture *render.Texture, w, h float32, invert bool) *Image {
	el := NewElement("Image", 0, 0, w, h)
	mat := assets.GetMaterial("ui_texture")
	mat.AddTexture("image", texture)
	return &Image{
		Element: el,
		Texture: texture,
		Quad:    geometry.NewImageQuad(mat, w, h, invert),
	}
}

func NewDepthImage(texture *render.Texture, w, h float32, invert bool) *Image {
	el := NewElement("DepthImage", 0, 0, w, h)
	mat := assets.GetMaterial("ui_depth_texture")
	mat.AddTexture("image", texture)
	return &Image{
		Element: el,
		Texture: texture,
		Quad:    geometry.NewImageQuad(mat, w, h, invert),
	}
}

func (r *Image) Draw(args render.DrawArgs) {
	args.Transform = r.Element.Transform.Matrix.Mul4(args.Transform) //args.Transform.Mul4(r.Element.Transform.Matrix)

	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)
	r.Quad.Draw(args)

	for _, el := range r.Element.children {
		el.Draw(args)
	}
}

func (r *Image) SetSize(w, h float32) {
	if w != r.width || h != r.height {
		r.Element.SetSize(w, h)
		r.Quad.SetSize(w, h)
		u := math.Min(w/float32(r.Texture.Width), 1)
		v := math.Min(h/float32(r.Texture.Height), 1)
		r.Quad.SetUV(u, v)
	}
}

func (r *Image) DesiredSize(w, h float32) (float32, float32) {
	return r.width, r.height
}
