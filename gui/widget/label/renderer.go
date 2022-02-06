package label

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/gui/quad"
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
	size   int
	bounds vec2.T

	font  font.T
	tex   texture.T
	mat   material.T
	mesh  quad.T
	uvs   quad.UV
	color color.T
}

func (r *renderer) Draw(args render.Args, label T, props *Props) {
	if props.Text == "" {
		return
	}

	if props.Font == nil {
		props.Font = assets.DefaultFont()
	}

	if r.mesh == nil {
		r.mat = assets.GetMaterial("ui_texture")
		r.uvs = quad.DefaultUVs
		r.mesh = quad.New(r.mat, quad.Props{
			UVs:   r.uvs,
			Size:  r.bounds.Scaled(0.5),
			Color: props.Color,
		})
	}

	textChanged := r.text != props.Text
	fontChanged := r.font != props.Font
	sizeChanged := r.size != props.Size
	colorChanged := r.color != props.Color

	invalidateTexture := textChanged || sizeChanged || fontChanged
	invalidateMesh := sizeChanged || fontChanged || colorChanged

	r.text = props.Text
	r.font = props.Font
	r.size = props.Size
	r.color = props.Color

	if invalidateTexture {

		// (re)create label texture
		args := font.Args{
			LineHeight: props.LineHeight,
			Color:      color.White,
		}
		r.bounds = r.font.Measure(props.Text, args)

		// create texture if required
		if r.tex == nil {
			r.tex = gl_texture.New(int(r.bounds.X), int(r.bounds.Y))
			// todo: single channel texture to save memory
		}

		img := r.font.Render(r.text, args)
		r.tex.BufferImage(img)
	}

	if invalidateMesh {
		r.mesh.Update(quad.Props{
			Size:  r.bounds.Scaled(0.5),
			UVs:   r.uvs,
			Color: r.color,
		})
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

	r.mesh.Material().Use()
	r.mesh.Material().RGBA("tint", props.Color)
	r.mesh.Material().Texture("image", r.tex)
	r.mesh.Draw(args)
}
