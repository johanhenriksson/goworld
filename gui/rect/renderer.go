package rect

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/geometry"
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
	mesh  *geometry.Rect
	mesh2 *geometry.Quad
	size  vec2.T
}

func (r *renderer) Draw(args render.Args, frame T, props *Props) {
	if r.mesh == nil {
		r.tex = assets.GetColorTexture(color.White)

		r.mat = assets.GetMaterial("ui_texture")
		r.mat.Texture("image", r.tex)

		r.mesh = geometry.NewRect(r.mat, vec2.Zero)
	}
	if r.mesh2 == nil {
		r.mesh2 = geometry.NewQuad(r.mat, vec2.Zero)
	}

	// set correct blending
	// perhaps this belongs somewhere else
	render.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// resize if needed
	if !frame.Size().ApproxEqual(r.size) {
		r.mesh.SetSize(frame.Size())
		r.mesh2.SetSize(frame.Size())
		r.size = frame.Size()
	}
	if props.Border != r.mesh.BorderWidth() {
		r.mesh.SetBorderWidth(props.Border)
	}

	r.mat.Use()
	r.mat.RGBA("tint", props.Color)
	r.mat.Texture("image", r.tex)
	//r.mesh.Draw(args)
	r.mesh2.Draw(args)
}

func (r *renderer) Destroy() {
	//  todo: clean up mesh, texture
}
