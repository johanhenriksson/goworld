package style

import "github.com/kjk/flex"

type Absolute struct {
	Left, Right PositionValueProp
	Top, Bottom PositionValueProp
}

func (a Absolute) ApplyPosition(fw FlexWidget) {
	fw.Flex().StyleSetPositionType(flex.PositionTypeAbsolute)

	if a.Left != nil {
		a.Left.ApplyPosition(fw, flex.EdgeLeft)
	}
	if a.Right != nil {
		a.Right.ApplyPosition(fw, flex.EdgeRight)
	}
	if a.Top != nil {
		a.Top.ApplyPosition(fw, flex.EdgeTop)
	}
	if a.Bottom != nil {
		a.Bottom.ApplyPosition(fw, flex.EdgeBottom)
	}
}

type Relative struct{}

func (r Relative) ApplyPosition(fw FlexWidget) {
	fw.Flex().StyleSetPositionType(flex.PositionTypeAbsolute)
}
