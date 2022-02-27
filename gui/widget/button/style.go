package button

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

var DefaultStyle = &Style{
	Padding: Rect{
		Left:   8,
		Right:  8,
		Top:    2,
		Bottom: 2,
	},

	TextAlign: AlignCenter,
	Font: Font{
		Name: "fonts/SourceCodeProRegular.ttf",
		Size: 14,
	},

	Color:      color.White,
	Background: color.Yellow,

	Hover: Hover{
		Background: color.Green,
	},

	Pressed: Pressed{
		Background: color.Purple,
	},
}

type Style struct {
	Extends *Style
	Hover   Hover
	Pressed Pressed

	Font       FontProp
	Color      ColorProp
	Background ColorProp
	Padding    PaddingProp
	TextAlign  AlignItemsProp
}

type Hover struct {
	Color      ColorProp
	Background ColorProp
}

type Pressed struct {
	Background ColorProp
}

func (s *Style) backgroundStyle() rect.Style {
	base := *rect.DefaultStyle
	if s.Extends == nil {
		if s != DefaultStyle {
			base = DefaultStyle.backgroundStyle()
		}
	} else {
		base = base.Extend(s.Extends.backgroundStyle())
	}
	return base.Extend(rect.Style{
		Padding:        s.Padding,
		AlignItems:     s.TextAlign,
		JustifyContent: JustifyCenter,

		Color: s.Background,
		Hover: rect.Hover{
			Color: s.Hover.Background,
		},
		Pressed: rect.Pressed{
			Color: s.Pressed.Background,
		},
	})
}

func (s *Style) labelStyle() label.Style {
	base := *label.DefaultStyle
	if s.Extends == nil {
		if s != DefaultStyle {
			base = DefaultStyle.labelStyle()
		}
	} else {
		base = base.Extend(s.Extends.labelStyle())
	}
	return label.Style{
		Font: s.Font,
		Hover: label.Hover{
			Color: s.Hover.Color,
		},
	}
}
