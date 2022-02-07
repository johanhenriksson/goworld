package style

import (
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/kjk/flex"
)

// Each type should define its own style struct!!!
// rect.Style etc
type Sheet struct {
	Color color.T

	// Sizing properties

	Width     WidthProp
	MaxWidth  MaxWidthProp
	Height    HeightProp
	MaxHeight MaxHeightProp
	Padding   PaddingProp
	Margin    MarginProp

	// Flex properties

	Basis  BasisProp
	Grow   FlexGrowProp
	Shrink FlexShrinkProp
	Layout FlexDirectionProp
}

func (style *Sheet) Apply(node *flex.Node) {
	node.StyleSetDisplay(flex.DisplayFlex)

	if style.Basis != nil {
		style.Basis.ApplyBasis(node)
	}
	if style.Width != nil {
		style.Width.ApplyWidth(node)
	}
	if style.MaxWidth != nil {
		style.MaxWidth.ApplyMaxWidth(node)
	}
	if style.Height != nil {
		style.Height.ApplyHeight(node)
	}
	if style.MaxHeight != nil {
		style.MaxHeight.ApplyMaxHeight(node)
	}
	if style.Padding != nil {
		style.Padding.ApplyPadding(node)
	}
	if style.Margin != nil {
		style.Margin.ApplyMargin(node)
	}
	if style.Layout != nil {
		style.Layout.ApplyFlexDirection(node)
	}
	if style.Grow != nil {
		style.Grow.ApplyFlexGrow(node)
	}
	if style.Shrink != nil {
		style.Shrink.ApplyFlexShrink(node)
	}
}
