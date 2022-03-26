package image

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/gui/quad"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
)

type Renderer interface {
	widget.Renderer[T]

	SetSize(vec2.T)
	SetImage(texture.T)
	SetInvert(bool)
	SetColor(color.T)
}

type renderer struct {
	tint   color.T
	invert bool
	tex    texture.T

	invalid bool
	size    vec2.T
	mat     material.T
	mesh    quad.T
	uvs     quad.UV
}

func NewRenderer() Renderer {
	return &renderer{
		tint:    color.White,
		tex:     assets.DefaultTexture(),
		uvs:     quad.DefaultUVs,
		invalid: true,
	}
}

func (r *renderer) SetSize(size vec2.T) {
	r.invalid = r.invalid || size != r.size
	r.size = size
}

func (r *renderer) SetImage(tex texture.T) {
	if tex == nil {
		tex = assets.DefaultTexture()
	}
	r.invalid = r.invalid || tex != r.tex
	r.tex = tex
}

func (r *renderer) SetColor(tint color.T) {
	if tint == color.None {
		tint = color.White
	}
	r.invalid = r.invalid || tint != r.tint
	r.tint = tint
}

func (r *renderer) SetInvert(invert bool) {
	r.invalid = r.invalid || invert != r.invert
	r.invert = invert
}

func (r *renderer) Draw(args widget.DrawArgs, image T) {
	if r.mesh == nil {
		r.mesh = quad.New(quad.Props{
			UVs:   r.uvs,
			Size:  r.size,
			Color: color.White,
		})
	}

	r.SetSize(image.Size())

	if r.invalid {
		uvs := r.uvs
		if r.invert {
			uvs = uvs.Inverted()
		}

		r.mesh.Update(quad.Props{
			UVs:   uvs,
			Size:  r.size,
			Color: r.tint,
		})
		r.invalid = false
	}

	// set correct blending
	// perhaps this belongs somewhere else
	render.BlendMultiply()

	// r.mesh.Draw(args)
}

func (r *renderer) Destroy() {

}
