package rect

import (
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type T interface {
	widget.T

	Children() []widget.T
}

type rect struct {
	widget.T
	props    *Props
	renderer Renderer
	children []widget.T
}

type Props struct {
	Border float32
	Color  render.Color
	Layout layout.T
	Width  dimension.T
	Height dimension.T
}

func New(key string, props *Props, children ...widget.T) T {
	// defaults
	if props.Layout == nil {
		props.Layout = layout.Column{}
	}
	if props.Width == nil {
		props.Width = dimension.Auto()
	}
	if props.Height == nil {
		props.Height = dimension.Auto()
	}

	f := &rect{
		T:        widget.New(key),
		props:    props,
		children: children,
		renderer: &renderer{},
	}

	f.Reflow()
	return f
}

func (f *rect) Draw(args render.Args) {
	f.renderer.Draw(args, f, f.props)

	for _, child := range f.children {
		// calculate child tranasform
		pos := vec3.Extend(child.Position(), -1)
		transform := mat4.Translate(pos)
		childArgs := args
		childArgs.Transform = transform.Mul(&args.Transform)
		childArgs.Position = pos

		// draw child
		child.Draw(childArgs)
	}
}

func (f *rect) Reflow() {
	f.props.Layout.Flow(f)

	// recursively layout children
	for _, child := range f.children {
		child.Reflow()
	}
}

func (f *rect) Resize(s vec2.T) {
	f.T.Resize(s)
	f.Reflow()
}

func (f *rect) Props() widget.Props {
	return f.props
}

func (f *rect) Children() []widget.T { return f.children }

func (f *rect) Width() dimension.T  { return f.props.Width }
func (f *rect) Height() dimension.T { return f.props.Height }
