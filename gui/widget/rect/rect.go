package rect

import (
	"github.com/kjk/flex"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
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
	OnMouseUp    mouse.Callback
	OnMouseDown  mouse.Callback
	OnMouseEnter mouse.Callback
	OnMouseExit  mouse.Callback
	OnMouseMove  mouse.Callback
	OnMouseDrag  mouse.Callback
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

func (f *rect) Draw(args widget.DrawArgs) {
	if f.props.Style.Hidden {
		return
	}
	f.T.Draw(args)
	f.Renderer.Draw(args, f)
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

		// buttons
		if e.Action() == mouse.Press {
			f.state.Pressed = true
			f.props.Style.Apply(f, f.state)

			if f.props.OnMouseDown != nil {
				f.props.OnMouseDown(e)
				e.Consume()
			}
		}
		if e.Action() == mouse.Release {
			f.state.Pressed = false
			f.props.Style.Apply(f, f.state)

			if f.props.OnMouseUp != nil {
				f.props.OnMouseUp(e)
				e.Consume()
			}
		}

		// move
		if e.Action() == mouse.Move {
			if f.state.Pressed && f.props.OnMouseDrag != nil {
				f.props.OnMouseDrag(e)
			} else {
				if f.props.OnMouseMove != nil {
					f.props.OnMouseMove(e)
				}
			}
		}
	} else {
		if f.state.Pressed {
			if e.Action() == mouse.Move && f.props.OnMouseDrag != nil {
				f.props.OnMouseDrag(e)
			}
			if e.Action() == mouse.Release {
				f.state.Pressed = false
				f.props.Style.Apply(f, f.state)
			}
		}

		// hover end
		if f.state.Hovered {
			f.state.Hovered = false
			f.props.Style.Apply(f, f.state)

			if f.props.OnMouseExit != nil {
				f.props.OnMouseExit(e)
			}
		}
	}
}
