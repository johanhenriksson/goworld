package style

import (
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/kjk/flex"
)

type Sheet struct {
	// Display properties

	Color ColorProp

	// Text properties

	Font       FontProp
	FontColor  ColorProp
	LineHeight LineHeightProp

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

func (style *Sheet) Apply(w any) {
	if fw, ok := w.(FlexWidget); ok {
		// always set display: flex
		fw.Flex().StyleSetDisplay(flex.DisplayFlex)

		if style.Basis != nil {
			style.Basis.ApplyBasis(fw)
		}
		if style.Width != nil {
			style.Width.ApplyWidth(fw)
		}
		if style.MaxWidth != nil {
			style.MaxWidth.ApplyMaxWidth(fw)
		}
		if style.Height != nil {
			style.Height.ApplyHeight(fw)
		}
		if style.MaxHeight != nil {
			style.MaxHeight.ApplyMaxHeight(fw)
		}
		if style.Padding != nil {
			style.Padding.ApplyPadding(fw)
		}
		if style.Margin != nil {
			style.Margin.ApplyMargin(fw)
		}
		if style.Layout != nil {
			style.Layout.ApplyFlexDirection(fw)
		}
		if style.Grow != nil {
			style.Grow.ApplyFlexGrow(fw)
		}
		if style.Shrink != nil {
			style.Shrink.ApplyFlexShrink(fw)
		}
	}

	if style.Color != nil {
		if cc, ok := w.(Colorizable); ok {
			rgba := style.Color.Vec4()
			cc.SetColor(color.RGBA(rgba.X, rgba.Y, rgba.Z, rgba.W))
		}
	}

	if fw, ok := w.(FontWidget); ok {
		if style.Font != nil {
			style.Font.ApplyFont(fw)
		}
		if style.FontColor != nil {
			rgba := style.FontColor.Vec4()
			fw.SetFontColor(color.RGBA(rgba.X, rgba.Y, rgba.Z, rgba.W))
		}
	}
}
