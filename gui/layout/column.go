package layout

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Column struct {
	Padding  float32
	Gutter   float32
	Relative bool
}

func (c Column) Flow(w Layoutable) {
	// column layouts share the height
	x := float32(c.Padding)
	y := float32(c.Padding)

	bounds := w.Size()
	items := w.Children()
	inner := bounds.Sub(vec2.New(2*c.Padding, 2*c.Padding))
	inner.Y -= c.Gutter * float32(len(items)-1)

	// calculate total desired height
	totalWeight := float32(0)
	for _, item := range items {
		totalWeight += item.Height().Resolve(inner.Y)
	}

	if !c.Relative {
		totalWeight = math.Max(totalWeight, inner.Y)
	}

	for _, item := range items {
		height := inner.Y * item.Height().Resolve(inner.Y) / totalWeight
		item.Resize(vec2.New(inner.X, height))
		item.Move(vec2.New(x, y))
		y += height + c.Gutter
	}
}
