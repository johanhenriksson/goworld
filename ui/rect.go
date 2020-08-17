package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/render"
)

type Rect struct {
	*Element
	Style  Style
	layout RectLayout
	quad   *geometry.Quad
}

type RectLayout func(Component, Style, float32, float32) (float32, float32)

func NewRect(style Style, children ...Component) *Rect {
	mat := assets.GetMaterial("ui_color")
	color := style.Color("background", render.Transparent)

	r := &Rect{
		Element: NewElement("Rect", 0, 0, 0, 0),
		quad:    geometry.NewQuad(mat, color, 0, 0),
		layout:  ColumnLayout,
		Style:   style,
	}

	layout := style.String("layout", "column")
	if layout == "row" {
		r.layout = RowLayout
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
	// TODO set color
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	r.quad.Draw(local)

	/* call parent - draw children etc */
	r.Element.Draw(args)
}

func (r *Rect) DesiredSize(aw, ah float32) (float32, float32) {
	return r.layout(r, r.Style, aw, ah)
}

func (r *Rect) SetSize(w, h float32) {
	if w != r.width || h != r.height {
		r.Element.SetSize(w, h)
		r.quad.SetSize(w, h)
	}
}
