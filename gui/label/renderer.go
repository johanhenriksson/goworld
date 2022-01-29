package label

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
)

type Renderer interface {
	Draw(render.Args, T, *Props)
}

type renderer struct {
	text   string
	color  color.T
	bounds vec2.T

	font *render.Font
	tex  *render.Texture
	mat  *render.Material
	mesh *geometry.Rect
}

func (r *renderer) Draw(args render.Args, label T, props *Props) {
	font := assets.GetFont("assets/fonts/SourceCodeProRegular.ttf", props.Size, 1.0)

	if r.text != props.Text || font != r.font || r.color != props.Color {
		// (re)create label texture
		r.bounds = r.font.Measure(props.Text)
		r.tex = render.CreateTexture(int(r.bounds.X), int(r.bounds.Y))
		r.font = font

		r.text = props.Text
		r.color = props.Color

		r.font.Render(r.tex, r.text, r.color)
	}

	if r.mesh == nil {
		r.mat = assets.GetMaterial("ui_texture")
		r.mat.Textures.Add("image", r.tex)
		r.mesh = geometry.NewRect(r.mat, vec2.Zero)
	}

	// set correct blending
	// perhaps this belongs somewhere else
	render.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// resize mesh if needed
	if !label.Size().ApproxEqual(label.Size()) {
		r.mesh.SetSize(label.Size())
	}

	// we can center the label on the mesh by modifying the uvs
	scale := label.Size().Div(r.bounds)
	

	r.mesh.Material.Use()
	r.mesh.Material.RGBA("tint", color.White)
	r.mesh.Material.Textures.Set("image", r.tex)
	r.mesh.Draw(args)
}
