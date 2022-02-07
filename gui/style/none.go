package style

import "github.com/kjk/flex"

// None implements all Prop interfaces but does nothing
// It's a cool idea but is it worth setting it everywhere?
type None struct{}

func (n None) ApplyBasis(node *flex.Node)     {}
func (n None) ApplyWidth(node *flex.Node)     {}
func (n None) ApplyMaxWidth(node *flex.Node)  {}
func (n None) ApplyHeight(node *flex.Node)    {}
func (n None) ApplyMaxHeight(node *flex.Node) {}
func (n None) ApplyPadding(node *flex.Node)   {}
func (n None) ApplyMargin(node *flex.Node)    {}
