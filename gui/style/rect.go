package style

import "github.com/kjk/flex"

type Rect struct {
	Left   float32
	Right  float32
	Top    float32
	Bottom float32
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
