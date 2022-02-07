package style

import "github.com/kjk/flex"

// Px is a pixel value
type Px float32

func (p Px) ApplyBasis(n *flex.Node)     { n.StyleSetFlexBasis(float32(p)) }
func (p Px) ApplyWidth(n *flex.Node)     { n.StyleSetWidth(float32(p)) }
func (p Px) ApplyMaxWidth(n *flex.Node)  { n.StyleSetMaxWidth(float32(p)) }
func (p Px) ApplyHeight(n *flex.Node)    { n.StyleSetHeight(float32(p)) }
func (p Px) ApplyMaxHeight(n *flex.Node) { n.StyleSetMaxHeight(float32(p)) }
func (p Px) ApplyPadding(n *flex.Node)   { n.StyleSetPadding(flex.EdgeAll, float32(p)) }
func (p Px) ApplyMargin(n *flex.Node)    { n.StyleSetMargin(flex.EdgeAll, float32(p)) }
