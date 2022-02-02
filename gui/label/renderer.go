package label

import (
	"fmt"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_texture"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/font"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Renderer interface {
	Draw(render.Args, T, *Props)
}

type renderer struct {
	text   string
	color  color.T
	bounds vec2.T
	size   vec2.T

	font font.T
	tex  texture.T
	mat  material.T
	mesh *geometry.Rect
}

func (r *renderer) Draw(args render.Args, label T, props *Props) {
	fnt := assets.GetFont("assets/fonts/SourceCodeProRegular.ttf", int(props.Size*2))

	if r.mesh == nil {
		r.mat = assets.GetMaterial("ui_texture")
		r.mesh = geometry.NewRect(r.mat, vec2.Zero)
	}

	if r.text != props.Text || fnt != r.font || r.color != props.Color {
		// (re)create label texture
		args := font.Args{
			LineHeight: props.LineHeight,
			Color:      props.Color,
		}
		r.font = fnt
		r.bounds = r.font.Measure(props.Text, args)
		r.tex = gl_texture.New(int(r.bounds.X), int(r.bounds.Y))

		r.text = props.Text
		r.color = props.Color

		fmt.Println("update label with text:", r.text, "bounds:", r.bounds)

		img := r.font.Render(r.text, args)
		img.SetRGBA(3, 3, color.Red.RGBA())
		r.tex.BufferImage(img)

		r.mesh.SetSize(r.bounds.Scaled(0.5))
	}

	// set correct blending
	// perhaps this belongs somewhere else
	render.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// resize mesh if needed
	// if !label.Size().ApproxEqual(r.size) {
	// 	fmt.Println("label size", label.Size())
	// 	r.mesh.SetSize(label.Size())
	// 	r.size = label.Size()
	// }

	// can the we use the gl viewport to clip anything out of bounds?

	// we can center the label on the mesh by modifying the uvs
	// scale := label.Size().Div(r.bounds)

	r.mesh.Material.Use()
	r.mesh.Material.RGBA("tint", color.White)
	r.mesh.Material.Texture("image", r.tex)
	r.mesh.Draw(args)
}
