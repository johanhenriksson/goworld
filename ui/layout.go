package ui

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec2"
)

func RowLayout(c Component, sz vec2.T) vec2.T {
	pad := c.GetStyle().Float("padding", 0)
	spacing := c.GetStyle().Float("spacing", 0)

	// ask each child for its desired size (in this case, width)
	//
	// simple case: desired size is less than the available space
	// divide any extra space according to the sizing mode of each child
	//  - grow: take available space, divided equally (default)
	//  - shrink: use as little space as possible (i.e. the requested size)
	//
	// annoying case: desired size is greater than the available space
	//
	// problems:
	//  - what if we want to constrain the maximum width?
	//    child A wants to grow, B wants to shrink.

	// can we just ignore dynamic horizontal sizing? its not a browser after all
	// two kinds of elements:
	// containers - variable width according to their largest child
	// controls - fixed width (labels, images, buttons etc)

	desired := vec2.T{}
	for _, child := range c.Children() {
		child.SetPosition(vec2.New(pad+desired.X, pad))
		childSize := child.Flow(vec2.T{
			X: sz.X - desired.X - 2*pad,
			Y: sz.Y - 2*pad,
		})
		desired.X += childSize.X + spacing
		desired.Y = math.Max(desired.Y, childSize.Y)
	}
	desired.X += 2*pad - spacing
	desired.Y += 2 * pad

	return c.Resize(desired)
}

func ColumnLayout(c Component, sz vec2.T) vec2.T {
	pad := c.GetStyle().Float("padding", 0)
	spacing := c.GetStyle().Float("spacing", 0)

	desired := vec2.T{}
	for _, child := range c.Children() {
		child.SetPosition(vec2.New(pad, pad+desired.Y))
		childSize := child.Flow(vec2.T{
			X: sz.X - 2*pad,
			Y: sz.Y - desired.Y - 2*pad,
		})
		desired.X = math.Max(desired.X, childSize.X)
		desired.Y += childSize.Y + spacing
	}
	desired.Y += 2*pad - spacing
	desired.X += 2 * pad

	return c.Resize(desired)
}

func FixedLayout(c Component, sz vec2.T) vec2.T {
	return vec2.T{
		X: c.Width(),
		Y: c.Height(),
	}
}
