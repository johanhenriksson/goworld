package layout

import (
	"fmt"

	"github.com/johanhenriksson/goworld/gui/dimension"
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
	// this is not recursive!
	// so it will never work?
	maxWidth := float32(0)
	fixedHeight := float32(0)
	totalWeight := float32(0)
	numAuto := 0
	for _, item := range items {
		w := item.Width().Resolve(inner.X)
		maxWidth = math.Max(maxWidth, w)

		h := item.Height().Resolve(inner.Y)
		if item.Height().Fixed() {
			fixedHeight += h
		} else if item.Height().Auto() {
			numAuto++
		} else {
			// percent
			totalWeight += h
		}
	}
	sharedHeight := math.Max(inner.Y-fixedHeight, 0)
	totalWeight = math.Max(sharedHeight, totalWeight)
	autoHeight := math.Max(inner.Y-totalWeight, 0)
	fmt.Println(w.Key(), "weight:", totalWeight, "auto:", autoHeight, "shared:", sharedHeight)

	x, y := c.Padding, c.Padding
	for _, item := range items {
		taken := vec2.Zero
		if item.Height().Fixed() {
			taken = item.Arrange(vec2.New(inner.X, inner.Y))
		} else if item.Height().Auto() {
			autow := sharedHeight * (autoHeight / float32(numAuto)) / (totalWeight + autoHeight)
			fmt.Println(item.Key(), "auto share:", autow)
			auto := item.Height().Resolve(autow)
			fmt.Println(item.Key(), "auto resolved:", auto)
			taken = item.Arrange(vec2.New(inner.X, auto))
			fmt.Println(item.Key(), "auto taken:", taken)
		} else {
			// percent
			share := sharedHeight * item.Height().Resolve(inner.Y) / totalWeight
			fmt.Println(item.Key(), "percent share:", share)
			taken = item.Arrange(vec2.New(inner.X, share))
			fmt.Println(item.Key(), "percent taken:", taken)
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

func (c Column) Width(basis dimension.T, children []widget.T, available float32) dimension.T {
	// if we have a non-auto value, return it
	if !basis.Auto() {
		return basis
	}

	// otherwise, grab the width of the largest child element
	var width dimension.T = dimension.Fixed(0)
	for _, child := range children {
		if child.Width().Resolve(available) > width.Resolve(available) {
			width = child.Width()
		}
	}
	return width
}

func (c Column) Height(basis dimension.T, children []widget.T, available float32) dimension.T {
	// if we have a non-auto value, return it
	if !basis.Auto() {
		return basis
	}

	// otherwise, sum the height of child elements
	var height float32
	for _, child := range children {
		height += child.Width().Resolve(available)
	}
	return dimension.Fixed(height)
}

type Dim struct {
	Value  float32
	Type   int // fixed, weight, auto
	Min    float32
	Max    float32
	Shrink bool
	Grow   bool
}
