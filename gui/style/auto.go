package style

type Auto struct{}

func (a Auto) ApplyWidth(fw FlexWidget)  { fw.Flex().StyleSetWidthAuto() }
func (a Auto) ApplyHeight(fw FlexWidget) { fw.Flex().StyleSetHeightAuto() }
