package label

import (
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/render/color"
)

type T interface {
	widget.T
}

type Props struct {
	Text  string
	Color color.T
	Size  float32
}

type label struct {
	widget.T

	renderer Renderer
	props    *Props
}

func NewLabel(key string, props *Props) T {
	if props.Size == 0 {
		props.Size = 12
	}

	return &label{
		T:        widget.New(key),
		props:    props,
		renderer: &renderer{},
	}
}

func (l *label) Props() widget.Props { return l.props }
