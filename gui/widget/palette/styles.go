package palette

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
)

var SwatchStyle = rect.Style{
	Grow:   Grow(0),
	Shrink: Shrink(0),
	Basis:  Pct(20),
	Height: Px(20),
}
