package rect

import (
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Layout interface {
	Flow(T, *Props)
}

type Column struct{}

func (c Column) Flow(rect T, props *Props) {
	// column layouts share the height
	x := float32(props.Padding)
	y := float32(props.Padding)

	bounds := rect.Size()
	items := rect.Children()
	inner := bounds.Sub(vec2.New(2*props.Padding, 2*props.Padding))
	inner.Y -= props.Gutter * float32(len(items)-1)

	itemSize := vec2.New(inner.X, inner.Y/float32(len(items)))
	for _, item := range items {
		item.Resize(itemSize)
		item.Move(vec2.New(x, y))
		y += itemSize.Y + props.Gutter
	}
}

type Row struct{}

func (r Row) Flow(rect T, props *Props) {
	// row layouts share the width
	x := float32(props.Padding)
	y := float32(props.Padding)
	bounds := rect.Size()
	items := rect.Children()
	inner := bounds.Sub(vec2.New(2*props.Padding, 2*props.Padding))
	inner.X -= props.Gutter * float32(len(items)-1)

	itemSize := vec2.New(inner.X/float32(len(items)), inner.Y)
	for _, item := range items {
		item.Resize(itemSize)
		item.Move(vec2.New(x, y))
		x += itemSize.X + props.Gutter
	}
}
