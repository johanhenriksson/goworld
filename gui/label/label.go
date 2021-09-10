package label

import (
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type T interface {
	widget.T
}

type TextMeasurer interface {
	Measure(string) vec2.T
}

type Props struct {
	Text string
	Size float32
}

type label struct {
	widget.T

	renderer Renderer
	props    *Props
}

func NewLabel(key string, props *Props) T {
	return &label{
		T:        widget.New(key),
		props:    props,
		renderer: &renderer{},
	}
}

func (l *label) Props() widget.Props { return l.props }
