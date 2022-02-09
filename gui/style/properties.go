package style

import (
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/kjk/flex"
)

type WidthProp interface{ ApplyWidth(n *flex.Node) }
type MaxWidthProp interface{ ApplyMaxWidth(n *flex.Node) }
type HeightProp interface{ ApplyHeight(n *flex.Node) }
type MaxHeightProp interface{ ApplyMaxHeight(n *flex.Node) }
type BasisProp interface{ ApplyBasis(n *flex.Node) }
type PaddingProp interface{ ApplyPadding(n *flex.Node) }
type MarginProp interface{ ApplyMargin(n *flex.Node) }
type FlexDirectionProp interface{ ApplyFlexDirection(n *flex.Node) }
type FlexGrowProp interface{ ApplyFlexGrow(n *flex.Node) }
type FlexShrinkProp interface{ ApplyFlexShrink(n *flex.Node) }

type FontProp interface{ ApplyFont(w widget.T) }
type FontColorProp interface{ ApplyFontColor(w widget.T) }
type LineHeightProp interface{ ApplyLineHeight(w widget.T) }
