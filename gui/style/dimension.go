package style

import "github.com/kjk/flex"

type Dim interface {
	SetWidth(*flex.Node)
	SetMaxWidth(*flex.Node)
	SetHeight(*flex.Node)
	SetMaxHeight(*flex.Node)
	SetBasis(*flex.Node)
}

type Fixed float32

func (f Fixed) SetBasis(n *flex.Node)     { n.StyleSetFlexBasis(float32(f)) }
func (f Fixed) SetWidth(n *flex.Node)     { n.StyleSetWidth(float32(f)) }
func (f Fixed) SetMaxWidth(n *flex.Node)  { n.StyleSetMaxWidth(float32(f)) }
func (f Fixed) SetHeight(n *flex.Node)    { n.StyleSetHeight(float32(f)) }
func (f Fixed) SetMaxHeight(n *flex.Node) { n.StyleSetMaxHeight(float32(f)) }

func Auto() Dim {
	return auto{}
}

type auto struct{}

func (a auto) SetWidth(n *flex.Node)     { n.StyleSetWidthAuto() }
func (a auto) SetMaxWidth(n *flex.Node)  {}
func (a auto) SetHeight(n *flex.Node)    { n.StyleSetHeightAuto() }
func (a auto) SetMaxHeight(n *flex.Node) {}
func (a auto) SetBasis(n *flex.Node)     {}

type Percent float32

func (p Percent) SetBasis(n *flex.Node)     { n.StyleSetFlexBasisPercent(float32(p)) }
func (p Percent) SetWidth(n *flex.Node)     { n.StyleSetWidthPercent(float32(p)) }
func (p Percent) SetMaxWidth(n *flex.Node)  { n.StyleSetMaxWidthPercent(float32(p)) }
func (p Percent) SetHeight(n *flex.Node)    { n.StyleSetHeightPercent(float32(p)) }
func (p Percent) SetMaxHeight(n *flex.Node) { n.StyleSetMaxHeightPercent(float32(p)) }
