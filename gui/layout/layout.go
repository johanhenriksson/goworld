package layout

import (
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type T interface {
	Flow(Layoutable)
}

type Layoutable interface {
	Size() vec2.T
	Children() []widget.T
}
