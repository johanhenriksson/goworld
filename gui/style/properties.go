package style

import (
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/font"
	"github.com/kjk/flex"
)

type FlexWidget interface {
	Flex() *flex.Node
}

type WidthProp interface{ ApplyWidth(fw FlexWidget) }
type MinWidthProp interface{ ApplyMinWidth(fw FlexWidget) }
type MaxWidthProp interface{ ApplyMaxWidth(fw FlexWidget) }
type HeightProp interface{ ApplyHeight(fw FlexWidget) }
type MaxHeightProp interface{ ApplyMaxHeight(fw FlexWidget) }
type MinHeightProp interface{ ApplyMinHeight(fw FlexWidget) }
type BasisProp interface{ ApplyBasis(fw FlexWidget) }
type PaddingProp interface{ ApplyPadding(fw FlexWidget) }
type MarginProp interface{ ApplyMargin(fw FlexWidget) }
type FlexDirectionProp interface{ ApplyFlexDirection(fw FlexWidget) }
type FlexGrowProp interface{ ApplyFlexGrow(fw FlexWidget) }
type FlexShrinkProp interface{ ApplyFlexShrink(fw FlexWidget) }
type AlignItemsProp interface{ ApplyAlignItems(fw FlexWidget) }
type AlignContentProp interface{ ApplyAlignContent(fw FlexWidget) }
type JustifyContentProp interface{ ApplyJustifyContent(fw FlexWidget) }

type PositionProp interface{ ApplyPosition(fw FlexWidget) }
type PositionValueProp interface {
	ApplyPosition(fw FlexWidget, edge flex.Edge)
}

type FontWidget interface {
	SetFont(font.T)
	SetFontSize(int)
	SetFontColor(color.T)
	SetLineHeight(float32)
}

type FontProp interface{ ApplyFont(fw FontWidget) }
type FontColorProp interface{ ApplyFontColor(fw FontWidget) }
type LineHeightProp interface{ ApplyLineHeight(fw FontWidget) }
