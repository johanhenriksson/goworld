package style

import (
	"github.com/kjk/flex"
)

type Column struct{}

func (c Column) ApplyFlexDirection(fw FlexWidget) {
	fw.Flex().StyleSetDisplay(flex.DisplayFlex)
	fw.Flex().StyleSetFlexDirection(flex.FlexDirectionColumn)
}

type Row struct{}

func (r Row) ApplyFlexDirection(fw FlexWidget) {
	fw.Flex().StyleSetDisplay(flex.DisplayFlex)
	fw.Flex().StyleSetFlexDirection(flex.FlexDirectionRow)
}

type Grow float32

func (g Grow) ApplyFlexGrow(fw FlexWidget) {
	fw.Flex().StyleSetDisplay(flex.DisplayFlex)
	fw.Flex().StyleSetFlexGrow(float32(g))
}

type Shrink float32

func (s Shrink) ApplyFlexShrink(fw FlexWidget) {
	fw.Flex().StyleSetDisplay(flex.DisplayFlex)
	fw.Flex().StyleSetFlexGrow(float32(s))
}

type Align flex.Align

const (
	AlignStart        = Align(flex.AlignFlexStart)
	AlignCenter       = Align(flex.AlignCenter)
	AlignEnd          = Align(flex.AlignFlexEnd)
	AlignStretch      = Align(flex.AlignStretch)
	AlignSpaceAround  = Align(flex.AlignSpaceAround)
	AlignSpaceBetween = Align(flex.AlignSpaceBetween)
)

func (a Align) ApplyAlignItems(fw FlexWidget)   { fw.Flex().StyleSetAlignItems(flex.Align(a)) }
func (a Align) ApplyAlignContent(fw FlexWidget) { fw.Flex().StyleSetAlignContent(flex.Align(a)) }

type Justify flex.Justify

const (
	JustifyStart        = Justify(flex.JustifyFlexStart)
	JustifyCenter       = Justify(flex.JustifyCenter)
	JustifyEnd          = Justify(flex.JustifyFlexEnd)
	JustifySpaceAround  = Justify(flex.JustifySpaceAround)
	JustifySpaceBetween = Justify(flex.JustifySpaceBetween)
)

func (j Justify) ApplyJustifyContent(fw FlexWidget) {
	fw.Flex().StyleSetJustifyContent(flex.Justify(j))
}
