package style

import (
	"github.com/kjk/flex"
)

// Px is a pixel value
type Px float32

func (p Px) ApplyBasis(fw FlexWidget)     { fw.Flex().StyleSetFlexBasis(float32(p)) }
func (p Px) ApplyWidth(fw FlexWidget)     { fw.Flex().StyleSetWidth(float32(p)) }
func (p Px) ApplyMaxWidth(fw FlexWidget)  { fw.Flex().StyleSetMaxWidth(float32(p)) }
func (p Px) ApplyHeight(fw FlexWidget)    { fw.Flex().StyleSetHeight(float32(p)) }
func (p Px) ApplyMaxHeight(fw FlexWidget) { fw.Flex().StyleSetMaxHeight(float32(p)) }
func (p Px) ApplyPadding(fw FlexWidget)   { fw.Flex().StyleSetPadding(flex.EdgeAll, float32(p)) }
func (p Px) ApplyMargin(fw FlexWidget)    { fw.Flex().StyleSetMargin(flex.EdgeAll, float32(p)) }

func (p Px) ApplyPosition(fw FlexWidget, edge flex.Edge) {
	fw.Flex().StyleSetPosition(edge, float32(p))
}

func (p Px) ApplyLineHeight(fw FontWidget) {
	fw.SetLineHeight(float32(p))
}
