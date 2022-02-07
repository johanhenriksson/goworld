package style

import "github.com/kjk/flex"

type Auto struct{}

func (a Auto) ApplyWidth(n *flex.Node)  { n.StyleSetWidthAuto() }
func (a Auto) ApplyHeight(n *flex.Node) { n.StyleSetHeightAuto() }
