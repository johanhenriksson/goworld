package style

import "github.com/kjk/flex"

type Rect struct {
	Left   float32
	Right  float32
	Top    float32
	Bottom float32
}

func (p Rect) ApplyPadding(node *flex.Node) {
	node.StyleSetPadding(flex.EdgeLeft, p.Left)
	node.StyleSetPadding(flex.EdgeRight, p.Right)
	node.StyleSetPadding(flex.EdgeTop, p.Top)
	node.StyleSetPadding(flex.EdgeBottom, p.Bottom)
}

func (p Rect) ApplyMargin(node *flex.Node) {
	node.StyleSetPadding(flex.EdgeLeft, p.Left)
	node.StyleSetPadding(flex.EdgeRight, p.Right)
	node.StyleSetPadding(flex.EdgeTop, p.Top)
	node.StyleSetPadding(flex.EdgeBottom, p.Bottom)
}
