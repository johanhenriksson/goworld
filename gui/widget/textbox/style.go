package textbox

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

var DefaultStyle = Style{
	Text: label.Style{
		Width: Pct(100),
		Color: color.Black,
	},
	Bg: rect.Style{
		Color:   color.White,
		Padding: RectXY(4, 2),
		Margin:  RectY(5),
	},
}
