package rect

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
	mat  *render.Material
	tex  *render.Texture
	mesh *geometry.Rect
	size vec2.T
}

func (r *renderer) Draw(args render.Args, frame T, props *Props) {
	if r.mesh == nil {
		r.tex = render.TextureFromColor(color.White)

		r.mat = assets.GetMaterial("ui_texture")
		r.mat.Textures.Add("image", r.tex)

		r.mesh = geometry.NewRect(r.mat, vec2.Zero)
	}

	// set correct blending
	// perhaps this belongs somewhere else
	render.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// resize if needed
	if !frame.Size().ApproxEqual(r.size) {
		r.mesh.SetSize(frame.Size())
	}
	if props.Border != r.mesh.BorderWidth() {
		r.mesh.SetBorderWidth(props.Border)
	}

	r.mesh.Material.Use()
	r.mesh.Material.RGBA("tint", props.Color)
	r.mesh.Material.Textures.Set("image", r.tex)
	r.mesh.Draw(args)
}

func (r *renderer) Destroy() {

}
