package ui

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Element struct {
	Style     Style
	Name      string
	Transform *Transform2D
	Size      vec2.T

	parent        Component
	children      []Component
	mouseHandlers []MouseHandler
}

func NewElement(name string, position, size vec2.T, style Style) *Element {
	e := &Element{
		Style:     style,
		Name:      name,
		Transform: CreateTransform2D(position, -1),
		Size:      size,

		children:      []Component{},
		mouseHandlers: []MouseHandler{},
	}
	return e
}

func (e *Element) ZIndex() float32 {
	// not sure how this is going to work yet
	// parents must be drawn underneath children (?)
	if e.parent != nil {
		return e.parent.ZIndex() + 1
	}
	return 0
}

// Parent peturns the parent element
func (e *Element) Parent() Component {
	return e.parent
}

// SetParent sets the parent element
func (e *Element) SetParent(parent Component) {
	// TODO detach from current parent?
	e.parent = parent
	e.Transform.Position = vec3.Extend(e.Transform.Position.XY(), e.ZIndex())
	e.Transform.Update(0)
}

// Children returns a list of child elements
func (e *Element) Children() []Component {
	return e.children
}

func (e *Element) Width() float32 {
	return e.Size.X
}

func (e *Element) Height() float32 {
	return e.Size.Y
}

func (e *Element) Resize(size vec2.T) vec2.T {
	e.Size = size
	return size
}

func (e *Element) Flow(available vec2.T) vec2.T {
	return available
}

func (e *Element) SetPosition(position vec2.T) {
	e.Transform.Position = vec3.Extend(position, e.Transform.Position.Z)
	e.Transform.Update(0)
}

// Attach a child to this element
func (e *Element) Attach(child Component) {
	e.children = append(e.children, child)
	// set parent?
}

// Draw this element and its children
func (e *Element) Draw(args engine.DrawArgs) {
	/* Multiply transform to args */
	args.Transform = e.Transform.Matrix.Mul(&args.Transform)
	for _, el := range e.children {
		el.Draw(args)
	}
}

// InBounds returns true of the given 2D position is wihtin the bounds of this element
func (e *Element) InBounds(pos vec2.T) bool {
	return pos.X >= 0 && pos.Y >= 0 &&
		pos.X <= e.Width() && pos.Y <= e.Height()
}

// HandleMouse attempts to handle a mouse event with this element
func (e *Element) HandleMouse(ev MouseEvent) bool {
	// transform the point into our local coordinate system
	invTransform := e.Transform.Matrix.Invert()
	projected := invTransform.TransformPoint(vec3.Extend(ev.Point, 0))
	ev.Point = projected.XY()

	// check if we're inside element bounds
	if !e.InBounds(ev.Point) {
		return false
	}

	// pass event to children
	for _, el := range e.children {
		handled := el.HandleMouse(ev)
		if handled {
			return true
		}
	}

	// execute local mouse handlers
	for _, callback := range e.mouseHandlers {
		callback(ev)
	}

	return len(e.mouseHandlers) > 0
}

// HandleInput is called when this element receives text input
func (e *Element) HandleInput(char rune) {}

// HandleKey is called when this element receives raw key events
func (e *Element) HandleKey(event KeyEvent) {}

// OnClick registers a mouse event handler
func (e *Element) OnClick(callback MouseHandler) {
	e.mouseHandlers = append(e.mouseHandlers, callback)
}

func (e *Element) Focus() {}
func (e *Element) Blur()  {}

func (e *Element) GetStyle() Style {
	return e.Style
}
