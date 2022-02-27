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
	style.Colorizable
}

type rect struct {
	widget.T
	Renderer

	props    Props
	children []widget.T
	state    style.State
}

type Props struct {
	Style        Style
	OnClick      mouse.Callback
	OnMouseEnter mouse.Callback
	OnMouseLeave mouse.Callback
	Children     []node.T
}

func New(key string, props Props) node.T {
	return node.Builtin(key, props, props.Children, Create)
}

func Create(key string, props Props) T {
	rect := &rect{
		T:        widget.New(key),
		Renderer: NewRenderer(),
	}
	rect.Update(props)
	return rect
}

func (f *rect) Draw(args render.Args) {
	if f.props.Style.Hidden {
		return
	}

	f.T.Draw(args)
	f.Renderer.Draw(args, f)

	for _, child := range f.children {
		// calculate child tranasform
		// try to fix the position to an actual pixel
		// pos := vec3.Extend(child.Position().Scaled(args.Viewport.Scale).Floor().Scaled(1/args.Viewport.Scale), -1)
		pos := vec3.Extend(child.Position(), args.Position.Z-1)
		transform := mat4.Translate(pos)
		childArgs := args
		childArgs.Transform = transform // .Mul(&args.Transform)
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
		child.Flex().Parent = f.Flex()
	}
	f.Flex().Children = nodes
}

//
// Lifecycle
//

func (f *rect) Props() any { return f.props }

func (f *rect) Update(p any) {
	new := p.(Props)
	styleChanged := new.Style != f.props.Style
	f.props = new

	if styleChanged {
		// apply new styles
		new.Style.Apply(f, f.state)
	}
}

func (f *rect) Destroy() {
	f.T.Destroy()
	f.Renderer.Destroy()

	for _, child := range f.children {
		child.Destroy()
	}
}

//
// Events
//

func (f *rect) MouseEvent(e mouse.Event) {
	// because children may have absolute positioning, we must pass the event to all of them.
	// children always have higher z index, so they have priority
	for _, frame := range f.children {
		if handler, ok := frame.(mouse.Handler); ok {
			handler.MouseEvent(e)
			if e.Handled() {
				e.Consume()
				return
			}
		}
	}

	target := e.Position().Sub(f.Position())
	size := f.Size()
	mouseover := target.X >= 0 && target.X < size.X && target.Y >= 0 && target.Y < size.Y

	if mouseover {
		// hover start
		if !f.state.Hovered {
			f.state.Hovered = true
			f.props.Style.Apply(f, f.state)

			if f.props.OnMouseEnter != nil {
				f.props.OnMouseEnter(e)
			}
		}

		// click
		if e.Action() == mouse.Press && f.props.OnClick != nil {
			f.props.OnClick(e)
			e.Consume()
		}
	} else {
		// hover end
		if f.state.Hovered {
			f.state.Hovered = false
			f.props.Style.Apply(f, f.state)

			if f.props.OnMouseLeave != nil {
				f.props.OnMouseLeave(e)
			}
		}
	}
}
