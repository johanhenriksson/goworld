package ui

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/geometry"
	"github.com/johanhenriksson/goworld/render"
)

type Rect struct {
	*Element
	Style Style
	quad  *geometry.Quad
}

func NewRect(x, y, w, h float32, style Style, children ...Component) *Rect {
	mat := assets.GetMaterial("ui_color")
	color := style.Color("background", render.Black)

	r := &Rect{
		Element: NewElement("Rect", x, y, w, h),
		quad:    geometry.NewQuad(mat, color, w, h),
		Style:   style,
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

	// column layout
	// width = min(style("width"), aw)
	// calculate desired height at current width
	dw := float32(0)
	dh := float32(0)
	for _, child := range r.children {
		cx, cy := float32(0), dh
		cdw, cdh := child.DesiredSize(aw, ah-dh)
		if cdw > dw {
			dw = cdw
		}
		child.SetPosition(cx, cy)
		dh += cdh
	}

	if dw > aw {
		dw = aw
	}

	r.SetSize(dw, dh)
	return dw, dh
}

func (r *Rect) SetSize(w, h float32) {
	if w != r.width || h != r.height {
		r.Element.SetSize(w, h)
		r.quad.SetSize(w, h)
	}
}
