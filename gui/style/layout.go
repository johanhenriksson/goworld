package style

import (
	"github.com/kjk/flex"
)

type Layout interface {
	Apply(*flex.Node)
}

type Column struct {
	Padding float32
}

func (c Column) Apply(node *flex.Node) {
	node.StyleSetDisplay(flex.DisplayFlex)
	node.StyleSetFlexDirection(flex.FlexDirectionColumn)
	node.StyleSetPadding(flex.EdgeAll, c.Padding)
}

type Row struct {
	Padding float32
}

func (r Row) Apply(node *flex.Node) {
	node.StyleSetDisplay(flex.DisplayFlex)
	node.StyleSetFlexDirection(flex.FlexDirectionRow)
	node.StyleSetPadding(flex.EdgeAll, r.Padding)
}

type Absolute struct {
}

func (a Absolute) Apply(node *flex.Node) {
	node.StyleSetPositionType(flex.PositionTypeAbsolute)
}
