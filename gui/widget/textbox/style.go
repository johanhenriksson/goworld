package textbox

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

var DefaultStyle = Style{
	Text: label.Style{
		Color: color.Black,
	},
	Bg: rect.Style{
		Color:   color.White,
		Padding: RectXY(4, 2),
		Basis:   Pct(100),
		Shrink:  Shrink(1),
		Grow:    Grow(1),
		Border: Border{
			Width: Px(1),
			Color: color.Black,
		},
	},
}
