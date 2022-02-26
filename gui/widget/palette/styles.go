package palette

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/label"
	"github.com/johanhenriksson/goworld/render/color"
)

var TitleStyle = label.Style{
	Color: color.White,
	Font: Font{
		Name: "fonts/SourceCodeProRegular.ttf",
		Size: 16,
	},

	// Hover: Hover{
	// 	FontColor: color.Red,
	// },
}
