package rect

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/gui/quad"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
)

type Renderer interface {
	Draw(render.Args, T)
	Destroy()

	SetColor(color.T)
}

type renderer struct {
	mat     material.T
	tex     texture.T
	mesh    quad.T
	size    vec2.T
	color   color.T
	uvs     quad.UV
	invalid bool
}

func NewRenderer() Renderer {
	return &renderer{
		invalid: true,
	}
}

func (r *renderer) SetSize(size vec2.T) {
	r.invalid = r.invalid || size != r.size
	r.size = size
}

func (r *renderer) SetColor(clr color.T) {
	r.invalid = r.invalid || clr != r.color
	r.color = clr
}

func (r *renderer) Draw(args render.Args, rect T) {
	// dont draw anything if its transparent anyway
	if r.color.A <= 0 {
		return
	}

	if r.mesh == nil {
		r.tex = assets.GetColorTexture(color.White)
		r.uvs = quad.DefaultUVs
		r.mat = assets.GetMaterial("ui_texture")
		r.mat.Texture("image", r.tex)
		r.mesh = quad.New(r.mat, quad.Props{
			UVs:   r.uvs,
			Size:  r.size,
			Color: r.color,
		})
	}

	r.SetSize(rect.Size())

	if r.invalid {
		r.mesh.Update(quad.Props{
			UVs:   r.uvs,
			Size:  r.size,
			Color: r.color,
		})
		r.invalid = false
	}

	// set correct blending
	render.BlendMultiply()

	// render.Scissor(frame.Position(), frame.Size())

	r.mat.Use()
	r.mat.Texture("image", r.tex)
	r.mesh.Draw(args)

	// render.ScissorDisable()
}

func (r *renderer) Destroy() {
	//  todo: clean up mesh, texture
}
