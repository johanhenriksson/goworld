package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/render"
)

type Rect struct {
	*Element
	layout RectLayout
	quad   *geometry.Quad
}

type RectLayout func(Component, Size) Size

func NewRect(style Style, children ...Component) *Rect {
	mat := assets.GetMaterial("ui_color")
	color := style.Color("background", render.Transparent)

	r := &Rect{
		Element: NewElement("Rect", 0, 0, 0, 0, style),
		quad:    geometry.NewQuad(mat, color, 0, 0),
		layout:  ColumnLayout,
	}

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
	local.Transform = r.Element.Transform.Matrix.Mul4(args.Transform)

	/* draw rect */
	// this belongs in the quad drawing code
	// avoid GL calls outside of the "core" packages render/engine/geometry
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	r.quad.Draw(local)

	/* call parent - draw children etc */
	r.Element.Draw(args)
}

func (r *Rect) Flow(available Size) Size {
	return r.layout(r, available)
}

func (r *Rect) Resize(size Size) Size {
	if size.Width != r.Width() || size.Height != r.Height() {
		r.Element.Resize(size)
		r.quad.SetSize(size.Width, size.Height)
	}
	return r.Size
}
