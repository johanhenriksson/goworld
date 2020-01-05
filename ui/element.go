package ui

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/render"
)

type Element struct {
	Name      string
	Transform *Transform2D

	width         float32
	height        float32
	z             float32
	parent        Component
	children      []Component
	mouseHandlers []MouseHandler
}

func (m *Manager) NewElement(name string, x, y, w, h, z float32) *Element {
	e := &Element{
		Name:      name,
		Transform: CreateTransform2D(x, y, z),

		width:         w,
		height:        h,
		children:      []Component{},
		mouseHandlers: []MouseHandler{},
	}
	return e
}

func (e *Element) ZIndex() float32 {
	// not sure how this is going to work yet
	// parents must be drawn underneath children (?)
	return e.z
}

// Parent peturns the parent element
func (e *Element) Parent() Component {
	return e.parent
}

// SetParent sets the parent element
func (e *Element) SetParent(parent Component) {
	// TODO detach from current parent?
	e.parent = parent
}

// Children returns a list of child elements
func (e *Element) Children() []Component {
	return e.children
}

func (e *Element) Width() float32 {
	return e.width
}

func (e *Element) Height() float32 {
	return e.height
}

// Attach a child to this element
func (e *Element) Attach(child Component) {
	e.children = append(e.children, child)
	// set parent?
}

// Draw this element and its children
func (e *Element) Draw(args render.DrawArgs) {
	/* Multiply transform to args */
	args.Transform = e.Transform.Matrix.Mul4(args.Transform)
	for _, el := range e.children {
		el.Draw(args)
	}
}

// InBounds returns true of the given 2D position is wihtin the bounds of this element
func (e *Element) InBounds(pos mgl.Vec2) bool {
	return pos.X() >= 0 && pos.Y() >= 0 &&
		pos.X() <= e.width && pos.Y() <= e.height
}

// HandleMouse attempts to handle a mouse event with this element
func (e *Element) HandleMouse(ev MouseEvent) bool {
	// transform the point into our local coordinate system
	projected := e.Transform.Matrix.Inv().Mul4x1(mgl.Vec4{ev.Point.X(), ev.Point.Y(), 0, 1})
	ev.Point = mgl.Vec2{projected.X(), projected.Y()}

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
