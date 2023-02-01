package style

import "github.com/kjk/flex"

type Rect struct {
	Left   float32
	Right  float32
	Top    float32
	Bottom float32
}

func RectAll(v float32) Rect {
	return Rect{
		Left:   v,
		Right:  v,
		Top:    v,
		Bottom: v,
	}
}

func RectXY(x, y float32) Rect {
	return Rect{
		Left:   x,
		Right:  x,
		Top:    y,
		Bottom: y,
	}
}

func RectX(v float32) Rect {
	return Rect{
		Left:  v,
		Right: v,
	}
}

func RectY(v float32) Rect {
	return Rect{
		Top:    v,
		Bottom: v,
	}
}

func (p Rect) ApplyPadding(fw FlexWidget) {
	fw.Flex().StyleSetPadding(flex.EdgeLeft, p.Left)
	fw.Flex().StyleSetPadding(flex.EdgeRight, p.Right)
	fw.Flex().StyleSetPadding(flex.EdgeTop, p.Top)
	fw.Flex().StyleSetPadding(flex.EdgeBottom, p.Bottom)
}

func (p Rect) ApplyMargin(fw FlexWidget) {
	fw.Flex().StyleSetMargin(flex.EdgeLeft, p.Left)
	fw.Flex().StyleSetMargin(flex.EdgeRight, p.Right)
	fw.Flex().StyleSetMargin(flex.EdgeTop, p.Top)
	fw.Flex().StyleSetMargin(flex.EdgeBottom, p.Bottom)
}
