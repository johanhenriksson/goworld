package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
)

type Rect struct {
	*Element
	layout RectLayout
	quad   *geometry.Quad
	tex    *render.Texture
}

type RectLayout func(Component, vec2.T) vec2.T

func NewRect(style Style, children ...Component) *Rect {
	mat := assets.GetMaterial("ui_texture")
	size := vec2.Zero
	position := vec2.Zero

	r := &Rect{
		Element: NewElement("Rect", position, size, style),
		quad:    geometry.NewQuad(mat, size),
		layout:  ColumnLayout,
		tex:     render.TextureFromColor(render.White),
	}
	mat.AddTexture("image", r.tex)

	layout := style.String("layout", "column")
	if layout == "row" {
		r.layout = RowLayout
	} else if layout == "fixed" {
		r.layout = FixedLayout
	}

	border := style.Float("radius", 0)
	r.quad.SetBorderWidth(border)

	for _, child := range children {
		r.Attach(child)
	}

	return r
}

func (r *Rect) Draw(args render.DrawArgs) {
	// this is sort of ugly. we dont really want to duplicate the transform
	// multiplication to every element. on the other hand, most elements
	// will need to apply the transform before they draw themselves

	/* compute local transform */
	local := args
	local.Transform = r.Element.Transform.Matrix.Mul(&args.Transform)

	/* draw rect */
	// this belongs in the quad drawing code
	// avoid GL calls outside of the "core" packages render/engine/geometry
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	color := r.Style.Color("color", render.Transparent)
	image := r.Style.Texture("image", r.tex)
	r.quad.Material.Use()
	r.quad.Material.RGBA("tint", color)
	r.quad.Material.SetTexture("image", image)
	r.quad.Draw(local)

	/* call parent - draw children etc */
	r.Element.Draw(args)
}

func (r *Rect) Flow(available vec2.T) vec2.T {
	return r.layout(r, available)
}

func (r *Rect) Resize(size vec2.T) vec2.T {
	if size.X != r.Width() || size.Y != r.Height() {
		r.Element.Resize(size)
		r.quad.SetSize(size)
	}
	return r.Size
}
