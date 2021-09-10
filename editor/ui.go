package editor

import (
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/ui"
)

var WindowColor = color.RGBA(0.15, 0.15, 0.15, 0.85)
var TextColor = color.RGBA(1, 1, 1, 1)

var WindowStyle = ui.Style{
	"color":   ui.Color(WindowColor),
	"radius":  ui.Float(3),
	"padding": ui.Float(5),
}
