package layout

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Row struct {
	Padding  float32
	Gutter   float32
	Relative bool
}

func (r Row) Flow(w Layoutable) {
	// row layouts share the width
	x := float32(r.Padding)
	y := float32(r.Padding)
	bounds := w.Size()
	items := w.Children()
	inner := bounds.Sub(vec2.New(2*r.Padding, 2*r.Padding))
	inner.X -= r.Gutter * float32(len(items)-1)

	// calculate total desired height
	totalWeight := float32(0)
	for _, item := range items {
		totalWeight += item.Width().Resolve(inner.X)
	}

	if !r.Relative {
		totalWeight = math.Max(totalWeight, inner.X)
	}

	for _, item := range items {
		width := inner.X * item.Width().Resolve(inner.X) / totalWeight
		item.Resize(vec2.New(width, inner.Y))
		item.Move(vec2.New(x, y))
		x += width + r.Gutter
	}
}
