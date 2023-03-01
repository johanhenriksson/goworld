package rect

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
)

type T interface {
	widget.T
	style.Colorizable
}

type rect struct {
	widget.T

	color    color.T
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

func Create(w widget.T, props Props) T {
	rect := &rect{
		T: w,
	}
	rect.Update(props)
	return rect
}

func (f *rect) Color() color.T     { return f.color }
func (f *rect) SetColor(c color.T) { f.color = c }

func (f *rect) ZOffset() int {
	return f.props.Style.ZOffset
}

func (f *rect) Children() []widget.T { return f.children }

// Updates the widget with a new set of children.
// The widget takes ownership of the passed child array.
func (f *rect) SetChildren(c []widget.T) {
	// reconcile flex children
	for i, child := range c {
		if i < len(f.Flex().Children) {
			// we have an existing child at position i
			if f.Flex().Children[i] == child.Flex() {
				// its identical, move on
				continue
			} else {
				// its different, replace it
				if child.Flex().Parent != nil {
					child.Flex().Parent.RemoveChild(child.Flex())
				}
				f.Flex().InsertChild(child.Flex(), i)
			}
		} else {
			// insert new child
			if child.Flex().Parent != nil {
				child.Flex().Parent.RemoveChild(child.Flex())
			}
			f.Flex().InsertChild(child.Flex(), i)
		}
	}

	// if the array of new children is shorter than the old one, remove any remaining children
	for i := len(c); i < len(f.children); i++ {
		f.Flex().RemoveChild(f.children[i].Flex())
	}

	// update child array
	f.children = c
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

	for _, child := range f.children {
		child.Destroy()
	}
}

//
// Events
//

func (f *rect) MouseEvent(e mouse.Event) {
	if f.props.Style.Hidden {
		return
	}

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

func (f *rect) Draw(args widget.DrawArgs, quads *widget.QuadBuffer) {
	tex := args.Textures.Fetch(color.White)
	if tex != nil && f.color.A > 0 {
		quads.Push(widget.Quad{
			Min:   args.Position.XY(),
			Max:   args.Position.XY().Add(f.Size()),
			MinUV: vec2.Zero,
			MaxUV: vec2.One,
			Color: [4]color.T{
				f.Color(),
				f.Color(),
				f.Color(),
				f.Color(),
			},
			ZIndex:   args.Position.Z,
			Radius:   5,
			Softness: 5,
			Texture:  uint32(tex.ID),
		})
	}

	// draw children
	for _, child := range f.Children() {
		// calculate child transform
		// try to fix the position to an actual pixel
		// pos := vec3.Extend(child.Position().Scaled(args.Viewport.Scale).Floor().Scaled(1/args.Viewport.Scale), -1)
		z := child.ZOffset()
		childArgs := args
		childArgs.Position = vec3.Extend(child.Position(), args.Position.Z+float32(1+z))

		// draw child
		child.Draw(childArgs, quads)
	}
}
