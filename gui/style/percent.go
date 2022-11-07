package style

import "github.com/kjk/flex"

// Pct is a percentage value
type Pct float32

func (p Pct) ApplyBasis(fw FlexWidget)     { fw.Flex().StyleSetFlexBasisPercent(float32(p)) }
func (p Pct) ApplyWidth(fw FlexWidget)     { fw.Flex().StyleSetWidthPercent(float32(p)) }
func (p Pct) ApplyMinWidth(fw FlexWidget)  { fw.Flex().StyleSetMinWidthPercent(float32(p)) }
func (p Pct) ApplyMaxWidth(fw FlexWidget)  { fw.Flex().StyleSetMaxWidthPercent(float32(p)) }
func (p Pct) ApplyHeight(fw FlexWidget)    { fw.Flex().StyleSetHeightPercent(float32(p)) }
func (p Pct) ApplyMaxHeight(fw FlexWidget) { fw.Flex().StyleSetMaxHeightPercent(float32(p)) }
func (p Pct) ApplyMinHeight(fw FlexWidget) { fw.Flex().StyleSetMinHeightPercent(float32(p)) }
func (p Pct) ApplyPadding(fw FlexWidget)   { fw.Flex().StyleSetPaddingPercent(flex.EdgeAll, float32(p)) }
func (p Pct) ApplyMargin(fw FlexWidget)    { fw.Flex().StyleSetMarginPercent(flex.EdgeAll, float32(p)) }

func (p Pct) ApplyPosition(fw FlexWidget, edge flex.Edge) {
	fw.Flex().StyleSetPositionPercent(edge, float32(p))
}
