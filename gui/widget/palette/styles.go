package palette

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
	"github.com/johanhenriksson/goworld/render/color"
)

var SwatchStyle = rect.Style{
	Grow:   Grow(0),
	Shrink: Shrink(1),
	Basis:  Pct(16),
	Height: Px(20),
	Margin: Px(2),
}

var TitleStyle = label.Style{
	Grow: Grow(1),

	Color: color.White,
	Font: Font{
		Name: "fonts/SourceCodeProRegular.ttf",
		Size: 16,
	},

	Hover: label.Hover{
		Color: color.Red,
	},
}
