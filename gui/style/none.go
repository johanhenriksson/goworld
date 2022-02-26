package style

// None implements all Prop interfaces but does nothing
// It's a cool idea but is it worth setting it everywhere?
type None struct{}

func (n None) ApplyBasis(fw FlexWidget)     {}
func (n None) ApplyWidth(fw FlexWidget)     {}
func (n None) ApplyMaxWidth(fw FlexWidget)  {}
func (n None) ApplyHeight(fw FlexWidget)    {}
func (n None) ApplyMaxHeight(fw FlexWidget) {}
func (n None) ApplyPadding(fw FlexWidget)   {}
func (n None) ApplyMargin(fw FlexWidget)    {}
