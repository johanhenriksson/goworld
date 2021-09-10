package layout

import (
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Row struct {
	Padding float32
	Gutter  float32
}

func (r Row) Flow(w Layoutable) {
	// row layouts share the width
	x := float32(r.Padding)
	y := float32(r.Padding)
	bounds := w.Size()
	items := w.Children()
	inner := bounds.Sub(vec2.New(2*r.Padding, 2*r.Padding))
	inner.X -= r.Gutter * float32(len(items)-1)

	itemSize := vec2.New(inner.X/float32(len(items)), inner.Y)
	for _, item := range items {
		item.Resize(itemSize)
		item.Move(vec2.New(x, y))
		x += itemSize.X + r.Gutter
	}
}
