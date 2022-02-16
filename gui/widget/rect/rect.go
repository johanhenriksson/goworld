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
	Renderer

	props    Props
	children []widget.T

	prevMouseTarget mouse.Handler
}

type Props struct {
	Style        style.Sheet
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
	f.T.Draw(args)
	f.Renderer.Draw(args, f)

	for _, child := range f.children {
		// calculate child tranasform
		// try to fix the position to an actual pixel
		pos := vec3.Extend(child.Position().Scaled(args.Viewport.Scale).Floor().Scaled(1/args.Viewport.Scale), -1)
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
	new := p.(Props)
	styleChanged := new.Style != f.props.Style
	f.props = new

	if styleChanged {
		new.Style.Apply(f)
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
	if e.Action() == mouse.Enter {
		// prop callback
		if f.props.OnMouseEnter != nil {
			f.props.OnMouseEnter(e)
		}
	}
	if e.Action() == mouse.Leave {
		// prop callback
		if f.props.OnMouseLeave != nil {
			f.props.OnMouseLeave(e)
		}

		if f.prevMouseTarget != nil {
			// pass it on if we have a current hover target
			f.prevMouseTarget.MouseEvent(e)
		}
		f.prevMouseTarget = nil
	}

	hit := false
	for _, frame := range f.children {
		if handler, ok := frame.(mouse.Handler); ok {
			ev := e.Project(frame.Position())
			target := ev.Position()
			size := frame.Size()
			if target.X < 0 || target.X > size.X || target.Y < 0 || target.Y > size.Y {
				// outside
				continue
			}

			// we hit something
			hit = true

			if f.prevMouseTarget != handler {
				if f.prevMouseTarget != nil {
					// exit!
					f.prevMouseTarget.MouseEvent(mouse.NewMouseLeaveEvent())
				}

				// mouse enter!
				handler.MouseEvent(mouse.NewMouseEnterEvent())
			}
			f.prevMouseTarget = handler

			handler.MouseEvent(ev)
			if ev.Handled() {
				e.Consume()
				return
			}
		}
	}

	if !hit && f.prevMouseTarget != nil {
		f.prevMouseTarget.MouseEvent(mouse.NewMouseLeaveEvent())
		f.prevMouseTarget = nil
	}

	if e.Action() == mouse.Press && f.props.OnClick != nil {
		f.props.OnClick(e)
		e.Consume()
	}
}
