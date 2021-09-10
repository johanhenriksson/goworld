package rect

import (
	"github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
)

type T interface {
	widget.T
}

type rect struct {
	widget.T
	props    *Props
	renderer Renderer
}

type Props struct {
	Border float32
	Color  render.Color
	Layout layout.T
}

func New(key string, props *Props, items ...widget.T) T {
	// defaults
	if props.Layout == nil {
		props.Layout = layout.Column{}
	}

	f := &rect{
		T:        widget.New(key, items...),
		props:    props,
		renderer: &renderer{},
	}

	f.Reflow()
	return f
}

func (f *rect) Draw(args render.Args) {
	f.renderer.Draw(args, f, f.props)
}

func (f *rect) Reflow() {
	f.props.Layout.Flow(f)
	f.T.Reflow()
}

func (f *rect) Resize(s vec2.T) {
	f.T.Resize(s)
	f.Reflow()
}

func (f *rect) Props() widget.Props {
	return f.props
}
