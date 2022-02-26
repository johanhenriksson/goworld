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

func (g Grow) ApplyFlexGrow(fw FlexWidget) { fw.Flex().StyleSetFlexGrow(float32(g)) }

type Shrink float32

func (s Shrink) ApplyFlexShrink(fw FlexWidget) { fw.Flex().StyleSetFlexGrow(float32(s)) }
