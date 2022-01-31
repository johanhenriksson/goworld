package label

import (
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
)

type T interface {
	widget.T
}

type Props struct {
	Text       string
	Color      color.T
	Size       float32
	LineHeight float32
}

type label struct {
	widget.T

	renderer Renderer
	props    *Props
}

func New(key string, props *Props) T {
	if props.Size == 0 {
		props.Size = 12
	}
	if props.LineHeight == 0 {
		props.LineHeight = 0
	}

	return &label{
		T:        widget.New(key),
		props:    props,
		renderer: &renderer{},
	}
}

func (l *label) Props() widget.Props { return l.props }

func (l *label) Draw(args render.Args) {
	l.T.Draw(args)

	l.renderer.Draw(args, l, l.props)
}

func (l *label) Height() dimension.T {
	return dimension.Fixed(l.props.Size * 1.3333 * 2)
}
