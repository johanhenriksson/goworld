package image

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
}

type renderer struct {
	size vec2.T
	tex  texture.T
	mat  material.T
	mesh quad.T
	uvs  quad.UV
	tint color.T
}

func (r *renderer) Draw(args render.Args, image T, props *Props) {
	if r.mesh == nil {
		r.mat = assets.GetMaterial("ui_texture")
		r.uvs = quad.DefaultUVs
		r.mesh = quad.New(r.mat, quad.Props{
			UVs:   r.uvs,
			Size:  r.size,
			Color: props.Tint,
		})
	}

	sizeChanged := !image.Size().ApproxEqual(r.size)
	colorChanged := props.Tint != r.tint
	invalidated := sizeChanged || colorChanged

	if invalidated {
		r.size = image.Size()
		r.tint = props.Tint

		uvs := r.uvs
		if props.Invert {
			uvs = uvs.Inverted()
		}

		r.mesh.Update(quad.Props{
			UVs:   uvs,
			Size:  r.size,
			Color: r.tint,
		})
	}

	// set correct blending
	// perhaps this belongs somewhere else
	render.BlendMultiply()

	r.mesh.Material().Use()
	r.mesh.Material().Texture("image", props.Image)
	r.mesh.Draw(args)
}
