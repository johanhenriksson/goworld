package layout

import (
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Column struct {
	Padding float32
	Gutter  float32
}

func (c Column) Flow(w Layoutable) {
	// column layouts share the height
	x := float32(c.Padding)
	y := float32(c.Padding)

	bounds := w.Size()
	items := w.Children()
	inner := bounds.Sub(vec2.New(2*c.Padding, 2*c.Padding))
	inner.Y -= c.Gutter * float32(len(items)-1)

	itemSize := vec2.New(inner.X, inner.Y/float32(len(items)))
	for _, item := range items {
		item.Resize(itemSize)
		item.Move(vec2.New(x, y))
		sz := item.Size()
		y += sz.Y + c.Gutter
	}
}
