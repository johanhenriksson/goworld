package rect

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/gui/node"
	"github.com/johanhenriksson/goworld/gui/style"
	"github.com/johanhenriksson/goworld/gui/widget"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/kjk/flex"
)

type T interface {
	widget.T

	style.Colorizable
	style.BorderWidget
	style.RadiusWidget
}

type rect struct {
	key      string
	flex     *flex.Node
	props    Props
	children []widget.T
	state    style.State

	radius      float32
	color       color.T
	borderColor color.T
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

func Create(key string, props Props) widget.T {
	node := flex.NewNodeWithConfig(flex.NewConfig())
	node.Context = key
	rect := &rect{
		key:  key,
		flex: node,
	}
	rect.Update(props)
	return rect
}

//
// Widget implementation
//

func (r *rect) Key() string      { return r.key }
func (r *rect) Flex() *flex.Node { return r.flex }
func (f *rect) Props() any       { return f.props }
func (r *rect) Position() vec2.T { return vec2.New(r.flex.LayoutGetLeft(), r.flex.LayoutGetTop()) }
func (r *rect) Size() vec2.T     { return vec2.New(r.flex.LayoutGetWidth(), r.flex.LayoutGetHeight()) }

func (f *rect) Children() []widget.T { return f.children }

// Updates the widget with a new set of children.
// The widget takes ownership of the passed child array.
func (f *rect) SetChildren(c []widget.T) {
	for i, new := range c {
		existing := false
		for j, old := range f.children {
			if old == nil {
				// already been reconciled
				continue
			}
			if new.Key() == old.Key() {
				c[i] = f.children[j] // re-use old element
				f.children[j] = nil  // mark this child as reused
				existing = true
				break
			}
		}
		if !existing {
			// no existing element found -> new element at position i
			f.flex.InsertChild(new.Flex(), i)
		}
	}

	// all elements in the old array that are non-nil do not existe in the new child array, and should be destroyed
	for _, old := range f.children {
		if old == nil {
			// child re-used, skip
			continue
		}
		// this element was deleted
		f.flex.RemoveChild(old.Flex())
	}

	// replace the child array
	f.children = c
}

func (f *rect) Update(p any) {
	new := p.(Props)
	styleChanged := new.Style != f.props.Style
	f.props = new

	if styleChanged {
		// apply new styles
		new.Style.Apply(f, f.state)
	}
}

//
// Styles
//

func (f *rect) SetColor(c color.T)       { f.color = c }
func (f *rect) SetBorderColor(c color.T) { f.borderColor = c }
func (f *rect) SetRadius(r float32)      { f.radius = r }

//
// Draw
//

func (f *rect) drawSelf(args widget.DrawArgs, quads *widget.QuadBuffer) {
	// if alpha is zero, skip drawing
	if f.color.A <= 0 {
		return
	}

	// fetch white texture. its always cached
	tex := args.Textures.Fetch(color.White)

	zindex := args.Position.Z + float32(f.props.Style.ZOffset)
	min := args.Position.XY().Add(f.Position())
	max := min.Add(f.Size())

	// todo: add style properties
	shadow := color.Black
	shadowSoftness := float32(0)
	shadowOffset := vec2.New(2, 2)

	// drop shadow
	if shadowSoftness > 0 {
		quads.Push(widget.Quad{
			Min:      min.Add(shadowOffset),
			Max:      max.Add(shadowOffset),
			Color:    [4]color.T{shadow, shadow, shadow, shadow},
			ZIndex:   zindex - 0.1,
			Radius:   f.radius,
			Softness: shadowSoftness,
			Texture:  uint32(tex.ID),
		})
	}

	// background
	quads.Push(widget.Quad{
		Min:     min,
		Max:     max,
		MinUV:   vec2.Zero,
		MaxUV:   vec2.One,
		Color:   [4]color.T{f.color, f.color, f.color, f.color},
		ZIndex:  zindex,
		Radius:  f.radius,
		Texture: uint32(tex.ID),
	})

	// border
	borderWidth := f.Flex().LayoutGetBorder(flex.EdgeTop)
	if borderWidth > 0 && f.borderColor.A > 0 {
		quads.Push(widget.Quad{
			Min:     min,
			Max:     max,
			Color:   [4]color.T{f.borderColor, f.borderColor, f.borderColor, f.borderColor},
			ZIndex:  zindex + 0.1,
			Radius:  f.radius,
			Border:  borderWidth,
			Texture: uint32(tex.ID),
		})
	}
}

func (f *rect) Draw(args widget.DrawArgs, quads *widget.QuadBuffer) {
	if f.props.Style.Hidden {
		return
	}

	f.drawSelf(args, quads)

	position := args.Position.XY().Add(f.Position())
	zindex := args.Position.Z + float32(f.props.Style.ZOffset)

	childArgs := args
	childArgs.Position = vec3.Extend(position, zindex+1)

	// draw children
	for _, child := range f.Children() {
		child.Draw(childArgs, quads)
	}
}

//
// Events
//

func (f *rect) Destroy() {
}

func (f *rect) MouseEvent(e mouse.Event) {
	if f.props.Style.Hidden {
		return
	}

	// because children may have absolute positioning, we must pass the event to all of them.
	// children always have higher z index, so they have priority
	// todo: due to negative Z offsets, this might not be true. compare the actual z value
	for _, frame := range f.children {
		if handler, ok := frame.(mouse.Handler); ok {
			handler.MouseEvent(e)
			if e.Handled() {
				e.Consume()
				return
			}
		}
	}

	absolutePos := widget.AbsolutePosition(f.flex)
	target := e.Position().Sub(absolutePos)
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
			}
		}
		if e.Action() == mouse.Release {
			f.state.Pressed = false
			f.props.Style.Apply(f, f.state)

			if f.props.OnMouseUp != nil {
				f.props.OnMouseUp(e)
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
