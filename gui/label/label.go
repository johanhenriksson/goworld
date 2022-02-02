package label

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
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
	OnClick    mouse.Callback
}

type label struct {
	widget.T
	props    *Props
	renderer Renderer
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

func (l *label) Size() vec2.T { return l.T.Size() }

func (l *label) Props() widget.Props       { return l.props }
func (l *label) Update(props widget.Props) { l.props = props.(*Props) }

func (l *label) Draw(args render.Args) {
	l.T.Draw(args)

	l.renderer.Draw(args, l, l.props)
}

func (l *label) Height() dimension.T {
	return dimension.Fixed(l.props.Size*1.3333 + 4)
}

func (l *label) Measure(available vec2.T) vec2.T {
	// what is the available space?
	// is that the total space divided by the number of items?
	return vec2.Zero
}

//
// Events
//

func (l *label) MouseEvent(e mouse.Event) {
	if e.Action() == mouse.Press && l.props.OnClick != nil {
		l.props.OnClick(e)
		e.Consume()
	}
}
