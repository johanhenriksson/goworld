package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/texture"
)

// Rect is a rectangle with support for borders & rounded corners.
// Acts as the basic building block for all UI elements.
type Rect struct {
	*Element
	layout RectLayout
	mesh   *geometry.Rect
	tex    texture.T
}

type RectLayout func(Component, vec2.T) vec2.T

func NewRect(style Style, children ...Component) *Rect {
	mat := assets.GetMaterialShared("ui_texture")
	size := vec2.Zero
	position := vec2.Zero

	r := &Rect{
		Element: NewElement("Rect", position, size, style),
		mesh:    geometry.NewRect(mat, size),
		layout:  ColumnLayout,
		tex:     assets.GetColorTexture(color.White),
	}
	mat.Textures.Add("image", r.tex)

	layout := style.String("layout", "column")
	if layout == "row" {
		r.layout = RowLayout
	} else if layout == "fixed" {
		r.layout = FixedLayout
	}

	border := style.Float("radius", 0)
	r.mesh.SetBorderWidth(border)

	for _, child := range children {
		r.Attach(child)
	}

	return r
}

func (r *Rect) Draw(args render.Args) {
	// this is sort of ugly. we dont really want to duplicate the transform
	// multiplication to every element. on the other hand, most elements
	// will need to apply the transform before they draw themselves

	/* compute local transform */
	local := args
	local.Transform = r.Element.Transform.Matrix.Mul(&args.Transform)

	/* draw rect */
	// this belongs in the quad drawing code
	// avoid GL calls outside of the "core" packages render/engine/geometry
	render.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	color := r.Style.Color("color", color.Transparent)
	image := r.Style.Texture("image", r.tex)
	r.mesh.Material.Use()
	r.mesh.Material.RGBA("tint", color)
	r.mesh.Material.Textures.Set("image", image)
	r.mesh.Draw(local)

	/* call parent - draw children etc */
	r.Element.Draw(args)
}

func (r *Rect) Flow(available vec2.T) vec2.T {
	return r.layout(r, available)
}

func (r *Rect) Resize(size vec2.T) vec2.T {
	if size.X != r.Width() || size.Y != r.Height() {
		r.Element.Resize(size)
		r.mesh.SetSize(size)
	}
	return r.Size
}
