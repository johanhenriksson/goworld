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
	Draw(render.Args, T, *Props)
	Destroy()
}

type renderer struct {
	mat   material.T
	tex   texture.T
	mesh  quad.T
	size  vec2.T
	color color.T
	uvs   quad.UV
}

func (r *renderer) Draw(args render.Args, frame T, props *Props) {
	// dont draw anything if its transparent anyway
	if frame.Style().Color.A <= 0 {
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

	// set correct blending
	render.BlendMultiply()

	// resize if needed
	sizeChanged := !frame.Size().ApproxEqual(r.size)
	colorChanged := frame.Style().Color != r.color
	invalidated := sizeChanged || colorChanged

	if invalidated {
		r.size = frame.Size()
		r.color = frame.Style().Color
		r.mesh.Update(quad.Props{
			UVs:   r.uvs,
			Size:  r.size,
			Color: r.color,
		})
	}

	// render.Scissor(frame.Position(), frame.Size())

	r.mat.Use()
	r.mat.Texture("image", r.tex)
	r.mesh.Draw(args)

	// render.ScissorDisable()
}

func (r *renderer) Destroy() {
	//  todo: clean up mesh, texture
}
