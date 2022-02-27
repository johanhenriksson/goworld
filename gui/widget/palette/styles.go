package palette

import (
	. "github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget/rect"
)

var SwatchStyle = rect.Style{
	Grow:   Grow(0),
	Shrink: Shrink(1),
	Basis:  Pct(16),
	Height: Px(20),
	Margin: Px(2),
}
