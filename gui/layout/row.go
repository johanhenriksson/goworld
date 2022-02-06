package layout

import (
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Row struct {
	Padding float32
	Gutter  float32
}

func (r Row) Arrange(w widget.T, space vec2.T) vec2.T {
	// Arrange in a best effort manner, using no more than the given space
	size := vec2.New(
		w.Width().Resolve(space.X),
		w.Height().Resolve(space.Y))

	// clamp to available space
	size = vec2.Min(size, space)

	x := r.Padding
	y := r.Padding

	items := w.Children()
	inner := size.Sub(vec2.New(2*r.Padding, 2*r.Padding))
	inner.X -= r.Gutter * float32(len(items)-1)

	// calculate the amout of fixed space and available shared space
	fixedWidth := float32(0)
	totalWeight := float32(0)
	for _, item := range items {
		w := item.Width().Resolve(inner.X)
		if item.Width().Fixed() {
			fixedWidth += w
		} else {
			totalWeight += w
		}
	}

	height := float32(0)
	sharedWidth := inner.X - fixedWidth
	for _, item := range items {
		taken := vec2.Zero
		if item.Width().Fixed() {
			taken = item.Arrange(vec2.New(inner.X, inner.Y))
		} else {
			share := inner.X * item.Width().Resolve(sharedWidth) / totalWeight
			taken = item.Arrange(vec2.New(share, inner.Y))
		}

		item.SetPosition(vec2.New(x, y))
		x += math.Max(taken.X, 0) + r.Gutter
		height = math.Max(height, taken.Y)
	}

	// add the bottom y padding and remove the gutter behind the last element:
	x += r.Padding - r.Gutter
	// add width padding on both sides
	height += 2 * r.Padding
	content := vec2.New(x, height)

	if w.Width().Auto() && content.X > 0 {
		size.X = content.X
	}
	if w.Height().Auto() && content.Y > 0 {
		size.Y = content.Y
	}
	return size
}

func (a Row) Width(basis dimension.T, children []widget.T, available float32) dimension.T {
	return basis
}

func (a Row) Height(basis dimension.T, children []widget.T, available float32) dimension.T {
	return basis
}
