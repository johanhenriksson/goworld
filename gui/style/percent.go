package style

import "github.com/kjk/flex"

// Pct is a percentage value
type Pct float32

func (p Pct) ApplyBasis(n *flex.Node)     { n.StyleSetFlexBasisPercent(float32(p)) }
func (p Pct) ApplyWidth(n *flex.Node)     { n.StyleSetWidthPercent(float32(p)) }
func (p Pct) ApplyMaxWidth(n *flex.Node)  { n.StyleSetMaxWidthPercent(float32(p)) }
func (p Pct) ApplyHeight(n *flex.Node)    { n.StyleSetHeightPercent(float32(p)) }
func (p Pct) ApplyMaxHeight(n *flex.Node) { n.StyleSetMaxHeightPercent(float32(p)) }
func (p Pct) ApplyPadding(n *flex.Node)   { n.StyleSetPaddingPercent(flex.EdgeAll, float32(p)) }
func (p Pct) ApplyMargin(n *flex.Node)    { n.StyleSetMarginPercent(flex.EdgeAll, float32(p)) }
