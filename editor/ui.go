package editor

import (
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/ui"
)

var WindowColor = render.Color4(0.15, 0.15, 0.15, 0.85)
var TextColor = render.Color4(1, 1, 1, 1)

var WindowStyle = ui.Style{
	"color":   ui.Color(WindowColor),
	"radius":  ui.Float(3),
	"padding": ui.Float(5),
}
