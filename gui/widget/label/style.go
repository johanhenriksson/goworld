package label

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/kjk/flex"
)

var DefaultFont = Font{
	Name: "fonts/SourceSansPro-Regular.ttf",
	Size: 16,
}
var DefaultColor = color.White
var DefaultLineHeight = Pct(100)

var DefaultStyle = &Style{
	Font:       DefaultFont,
	LineHeight: DefaultLineHeight,
	Color:      DefaultColor,
}

type Style struct {
	Extends *Style
	Hover   Hover

	Hidden     bool
	Font       FontProp
	Color      ColorProp
	LineHeight LineHeightProp

	Width     WidthProp
	MaxWidth  MaxWidthProp
	Height    HeightProp
	MaxHeight MaxHeightProp

	Basis  BasisProp
	Grow   FlexGrowProp
	Shrink FlexShrinkProp
}

type Hover struct {
	Color ColorProp
}

func (s *Style) Extend(e Style) Style {
	e.Extends = s
	return e
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

	if style.Font != nil {
		style.Font.ApplyFont(w)
	}
	if style.Color != nil {
		rgba := style.Color.Vec4()
		w.SetFontColor(RGBA(rgba.X, rgba.Y, rgba.Z, rgba.W))
	}

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

	if state.Hovered {
		style.Hover.Apply(w)
	}
}

func (style *Hover) Apply(w T) {
	if style.Color != nil {
		rgba := style.Color.Vec4()
		w.SetFontColor(RGBA(rgba.X, rgba.Y, rgba.Z, rgba.W))
	}
}
