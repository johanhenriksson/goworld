package ui

import "github.com/johanhenriksson/goworld/math"

func RowLayout(c Component, sz Size) Size {
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

	desired := Size{}
	for _, child := range c.Children() {
		child.SetPosition(pad+desired.Width, pad)
		childSize := child.Flow(Size{
			Width:  sz.Width - desired.Width - 2*pad,
			Height: sz.Height - 2*pad,
		})
		desired.Width += childSize.Width + spacing
		desired.Height = math.Max(desired.Height, childSize.Height)
	}
	desired.Width += 2*pad - spacing
	desired.Height += 2 * pad

	return c.Resize(desired)
}

func ColumnLayout(c Component, sz Size) Size {
	pad := c.GetStyle().Float("padding", 0)
	spacing := c.GetStyle().Float("spacing", 0)

	desired := Size{}
	for _, child := range c.Children() {
		child.SetPosition(pad, pad+desired.Height)
		childSize := child.Flow(Size{
			Width:  sz.Width - 2*pad,
			Height: sz.Height - desired.Height - 2*pad,
		})
		desired.Width = math.Max(desired.Width, childSize.Width)
		desired.Height += childSize.Height + spacing
	}
	desired.Height += 2*pad - spacing
	desired.Width += 2 * pad

	return c.Resize(desired)
}

func FixedLayout(c Component, sz Size) Size {
	return Size{c.Width(), c.Height()}
}

type Size struct {
	Width  float32
	Height float32
}
