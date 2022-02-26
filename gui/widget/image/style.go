package image

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/kjk/flex"
)

var DefaultStyle = &Style{}

type Style struct {
	Extends *Style

	// Display properties

	Hidden bool

	// Sizing properties

	Width     WidthProp
	MaxWidth  MaxWidthProp
	Height    HeightProp
	MaxHeight MaxHeightProp

	// Flex properties

	Basis  BasisProp
	Grow   FlexGrowProp
	Shrink FlexShrinkProp
}

func (style *Style) Apply(w T, state State) {
	if style.Extends == nil {
		if style != DefaultStyle {
			DefaultStyle.Apply(w, state)
		}
	} else {
		style.Extends.Apply(w, state)
	}

	// always set display: flex
	w.Flex().StyleSetDisplay(flex.DisplayFlex)

	if style.Basis != nil {
		style.Basis.ApplyBasis(w)
	}
	if style.Width != nil {
		style.Width.ApplyWidth(w)
	}
	if style.MaxWidth != nil {
		style.MaxWidth.ApplyMaxWidth(w)
	}
	if style.Height != nil {
		style.Height.ApplyHeight(w)
	}
	if style.MaxHeight != nil {
		style.MaxHeight.ApplyMaxHeight(w)
	}
	if style.Grow != nil {
		style.Grow.ApplyFlexGrow(w)
	}
	if style.Shrink != nil {
		style.Shrink.ApplyFlexShrink(w)
	}
}
