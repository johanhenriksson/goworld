package layout

import (
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Column struct {
	Padding float32
	Gutter  float32
}

func (c Column) Arrange(w widget.T, space vec2.T) vec2.T {
	// Arrange in a best effort manner, using no more than the given space
	size := vec2.New(
		w.Width().Resolve(space.X),
		w.Height().Resolve(space.Y))

	// clamp to available space
	size = vec2.Min(size, space)

	items := w.Children()
	inner := size.Sub(vec2.New(2*c.Padding, 2*c.Padding))
	inner.Y -= c.Gutter * float32(len(items)-1)

	// calculate the amout of fixed space and available shared space
	maxWidth := float32(0)
	fixedHeight := float32(0)
	totalWeight := float32(0)
	// numAuto := 0
	for _, item := range items {
		w := item.Width().Resolve(inner.X)
		maxWidth = math.Max(maxWidth, w)

		h := item.Height().Resolve(inner.Y)
		if item.Height().Fixed() {
			fixedHeight += h
			// } else if item.Height().Auto() {
			// 	numAuto++
		} else {
			// percent
			totalWeight += h
		}
	}

	x, y := c.Padding, c.Padding
	// autoHeight := (1 - totalWeight) / float32(numAuto)
	sharedHeight := math.Max(inner.Y-fixedHeight, 0)
	for _, item := range items {
		taken := vec2.Zero
		if item.Height().Fixed() {
			taken = item.Arrange(vec2.New(inner.X, inner.Y))
			// } else if item.Height().Auto() {
			// 	auto := item.Height().Resolve(autoHeight * inner.Y)
			// 	taken = item.Arrange(vec2.New(inner.X, auto))
		} else {
			// percent
			share := sharedHeight * item.Height().Resolve(inner.Y) / totalWeight
			taken = item.Arrange(vec2.New(inner.X, share))
		}

		item.SetPosition(vec2.New(x, y))
		y += math.Max(taken.Y, 0) + c.Gutter
	}

	// final content size:
	// add the bottom y padding and remove the gutter behind the last element:
	// add width padding on both sides
	content := vec2.New(
		maxWidth+2*c.Padding,
		y+c.Padding-c.Gutter)

	// shrink/expand to content if auto sizing (if not empty)
	if w.Width().Auto() && content.X > 0 {
		size.X = content.X
	}
	if w.Height().Auto() && content.Y > 0 {
		size.Y = content.Y
	}

	return size
}
