package button

import (
	. "github.com/johanhenriksson/goworld/gui/style"
)

type Style struct {
	Hover     Hover
	TextColor ColorProp
	BgColor   ColorProp
	Padding   PaddingProp
	Margin    MarginProp
	Border    BorderProp
	Radius    RadiusProp
}

type Hover struct {
	TextColor ColorProp
	BgColor   ColorProp
}
