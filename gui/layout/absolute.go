package layout

import (
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Absolute struct{}

func (a Absolute) Arrange(w widget.T, space vec2.T) vec2.T {
	size := vec2.New(
		w.Width().Resolve(space.X),
		w.Height().Resolve(space.Y))

	for _, item := range w.Children() {
		pos := item.Position()
		bounds := size.Sub(pos)
		item.Arrange(bounds)
	}

	return size
}

func (a Absolute) Width(basis dimension.T, children []widget.T, available float32) dimension.T {
	return basis
}

func (a Absolute) Height(basis dimension.T, children []widget.T, available float32) dimension.T {
	return basis
}
