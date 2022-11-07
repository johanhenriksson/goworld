package rect

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/kjk/flex"
)

var DefaultStyle = &Style{
	Padding: None{},
	Margin:  None{},
	// Position: Relative{},
}

type Style struct {
	Extends *Style
	Hover   Hover
	Pressed Pressed

	// Display properties

	ZOffset  int
	Hidden   bool
	Color    ColorProp
	Position PositionProp

	// Sizing properties

	Width     WidthProp
	MinWidth  MinWidthProp
	MaxWidth  MaxWidthProp
	Height    HeightProp
	MinHeight MinHeightProp
	MaxHeight MaxHeightProp
	Padding   PaddingProp
	Margin    MarginProp

	// Flex properties

	Basis          BasisProp
	Grow           FlexGrowProp
	Shrink         FlexShrinkProp
	Layout         FlexDirectionProp
	AlignItems     AlignItemsProp
	AlignContent   AlignContentProp
	JustifyContent JustifyContentProp
}

type Hover struct {
	Color ColorProp
}

type Pressed struct {
	Color ColorProp
}

func (style *Style) Apply(w T, state State) {
	// this causes unnecessary updates
	// base styles are applied and then immediately overwritten
	if style.Extends == nil {
		if style != DefaultStyle {
			DefaultStyle.Apply(w, state)
		}
	} else {
		style.Extends.Apply(w, state)
	}

	if style.Position != nil {
		style.Position.ApplyPosition(w)
	}

	// always set display: flex
	w.Flex().StyleSetDisplay(flex.DisplayFlex)

	if style.Basis != nil {
		style.Basis.ApplyBasis(w)
	}
	if style.Width != nil {
		style.Width.ApplyWidth(w)
	}
	if style.MinWidth != nil {
		style.MinWidth.ApplyMinWidth(w)
	}
	if style.MaxWidth != nil {
		style.MaxWidth.ApplyMaxWidth(w)
	}
	if style.Height != nil {
		style.Height.ApplyHeight(w)
	}
	if style.MinHeight != nil {
		style.MinHeight.ApplyMinHeight(w)
	}
	if style.MaxHeight != nil {
		style.MaxHeight.ApplyMaxHeight(w)
	}
	if style.Padding != nil {
		style.Padding.ApplyPadding(w)
	}
	if style.Margin != nil {
		style.Margin.ApplyMargin(w)
	}
	if style.Layout != nil {
		style.Layout.ApplyFlexDirection(w)
	}
	if style.Grow != nil {
		style.Grow.ApplyFlexGrow(w)
	}
	if style.Shrink != nil {
		style.Shrink.ApplyFlexShrink(w)
	}
	if style.AlignItems != nil {
		style.AlignItems.ApplyAlignItems(w)
	}
	if style.AlignContent != nil {
		style.AlignContent.ApplyAlignContent(w)
	}
	if style.JustifyContent != nil {
		style.JustifyContent.ApplyJustifyContent(w)
	}

	if style.Color != nil {
		rgba := style.Color.Vec4()
		w.SetColor(RGBA(rgba.X, rgba.Y, rgba.Z, rgba.W))
	}

	if state.Hovered {
		style.Hover.Apply(w)
	}

	if state.Pressed {
		style.Pressed.Apply(w)
	}
}

func (s *Style) Extend(e Style) Style {
	e.Extends = s
	return e
}

func (s Hover) Apply(w T) {
	if s.Color != nil {
		rgba := s.Color.Vec4()
		w.SetColor(RGBA(rgba.X, rgba.Y, rgba.Z, rgba.W))
	}
}

func (s Pressed) Apply(w T) {
	if s.Color != nil {
		rgba := s.Color.Vec4()
		w.SetColor(RGBA(rgba.X, rgba.Y, rgba.Z, rgba.W))
	}
}
