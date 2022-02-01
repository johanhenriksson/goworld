package rect

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/gui/layout"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
)

type T interface {
	widget.T

	Children() []widget.T
	SetChildren([]widget.T)
}

type rect struct {
	widget.T
	props    *Props
	renderer Renderer
	children []widget.T
}

type Props struct {
	Border float32
	Color  color.T
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
	f.T.Draw(args)

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

func (f *rect) Children() []widget.T     { return f.children }
func (f *rect) SetChildren(c []widget.T) { f.children = c }

func (f *rect) Width() dimension.T  { return f.props.Width }
func (f *rect) Height() dimension.T { return f.props.Height }

func (f *rect) DesiredHeight(width float32) float32 {
	height := float32(0)
	for _, c := range f.children {
		height += c.DesiredHeight(width)
	}
	return height
}

//
// Lifecycle
//

func (f *rect) Props() widget.Props   { return f.props }
func (f *rect) Update(p widget.Props) { f.props = p.(*Props) }

func (f *rect) Destroy() {
	f.T.Destroy()
	f.renderer.Destroy()

	for _, child := range f.children {
		child.Destroy()
	}
}

//
// Events
//

func (f *rect) MouseEvent(e mouse.Event) {
	for _, frame := range f.children {
		if handler, ok := frame.(mouse.Handler); ok {
			ev := e.Project(frame.Position())
			target := ev.Position()
			size := frame.Size()
			if target.X < 0 || target.X > size.X || target.Y < 0 || target.Y > size.Y {
				// outside
				continue
			}

			handler.MouseEvent(ev)
			if ev.Handled() {
				e.StopPropagation()
				break
			}
		}
	}
}
