package ui

import "github.com/johanhenriksson/goworld/math"

func RowLayout(c Component, sz Size) Size {
	pad := c.Float("padding", 0)
	spacing := c.Float("spacing", 0)

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
	pad := c.Float("padding", 0)
	spacing := c.Float("spacing", 0)

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
