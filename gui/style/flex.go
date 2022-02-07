package style

import (
	"github.com/kjk/flex"
)

type Column struct{}

func (c Column) ApplyFlexDirection(node *flex.Node) {
	node.StyleSetDisplay(flex.DisplayFlex)
	node.StyleSetFlexDirection(flex.FlexDirectionColumn)
}

type Row struct{}

func (r Row) ApplyFlexDirection(node *flex.Node) {
	node.StyleSetDisplay(flex.DisplayFlex)
	node.StyleSetFlexDirection(flex.FlexDirectionRow)
}

type Grow float32

func (g Grow) ApplyFlexGrow(n *flex.Node) { n.StyleSetFlexGrow(float32(g)) }

type Shrink float32

func (s Shrink) ApplyFlexShrink(n *flex.Node) { n.StyleSetFlexGrow(float32(s)) }
