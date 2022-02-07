package rect

import (
	"github.com/kjk/flex"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type T interface {
	widget.T
}

type rect struct {
	widget.T
	props    *Props
	renderer Renderer
	children []widget.T
}

type Props struct {
	Style    style.Sheet
	OnClick  mouse.Callback
	Children []node.T
}

func New(key string, props *Props) node.T {
	return node.Builtin(key, props, props.Children, Create)
}

func Create(key string, props *Props) T {
	rect := &rect{
		T:        widget.New(key),
		renderer: &renderer{},
		props:    nil,
	}
	rect.Update(props)
	return rect
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

func (f *rect) Children() []widget.T { return f.children }
func (f *rect) SetChildren(c []widget.T) {
	f.children = c
	nodes := make([]*flex.Node, len(c))
	for i, child := range c {
		nodes[i] = child.Flex()
	}
	f.Flex().Children = nodes
}

//
// Lifecycle
//

func (f *rect) Props() any { return f.props }

func (f *rect) Update(p any) {
	new := p.(*Props)

	styleChanged := true
	if f.props != nil {
		styleChanged = new.Style != f.props.Style
	}

	// update props
	f.props = new

	if styleChanged {
		f.SetStyle(new.Style)
	}
}

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
				e.Consume()
				return
			}
		}
	}

	// how to do mouse enter/exit events?
	// we wont get any event when the mouse is outside

	if e.Action() == mouse.Press && f.props.OnClick != nil {
		f.props.OnClick(e)
		e.Consume()
	}
}
