package window

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/button"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

type Style struct {
	MinWidth MinWidthProp
	MaxWidth MaxWidthProp
}

var TitleStyle = label.Style{
	Grow: Grow(1),

	Color: RGB(1, 1, 1),
	Font: Font{
		Name: "fonts/SourceCodeProRegular.ttf",
		Size: 16,
	},
}

var TitlebarStyle = rect.Style{
	Color:      RGBA(0, 0, 0, 0.8),
	Padding:    Px(4),
	Layout:     Row{},
	AlignItems: AlignCenter,
	Pressed: rect.Pressed{
		Color: RGBA(0.2, 0.2, 0.2, 0.8),
	},
}

var FrameStyle = rect.Style{
	Color:        RGBA(0.1, 0.1, 0.1, 0.8),
	Padding:      RectXY(10, 10),
	Layout:       Column{},
	AlignItems:   AlignStart,
	AlignContent: AlignStretch,
}

var CloseButtonStyle = button.Style{
	Text: label.Style{
		Color: color.White,
	},
	Bg: rect.Style{
		Color: RGB(0.597, 0.098, 0.117),
		Padding: Rect{
			Left:   5,
			Right:  5,
			Top:    2,
			Bottom: 2,
		},

		Hover: rect.Hover{
			Color: RGB(0.3, 0.3, 0.3),
		},
	},
}
