package style

import (
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/kjk/flex"
)

type Sheet struct {
	Color     color.T
	Basis     Dim
	Width     Dim
	MaxWidth  Dim
	Height    Dim
	MaxHeight Dim
	Grow      int
	Shrink    int
	Layout    Layout
}

func (style *Sheet) Apply(node *flex.Node) {
	node.StyleSetDisplay(flex.DisplayFlex)
	node.StyleSetFlexGrow(float32(style.Grow))
	node.StyleSetFlexShrink(float32(style.Shrink))

	if style.Basis != nil {
		style.Basis.SetBasis(node)
	}
	if style.Width != nil {
		style.Width.SetWidth(node)
	}
	if style.MaxWidth != nil {
		style.MaxWidth.SetMaxWidth(node)
	}
	if style.Height != nil {
		style.Height.SetHeight(node)
	}
	if style.MaxHeight != nil {
		style.MaxHeight.SetMaxHeight(node)
	}

	if style.Layout != nil {
		style.Layout.Apply(node)
	}
}
