package ui

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/render"
)

type Element struct {
	width     float32
	height    float32
	z         float32
	parent    render.Drawable
	children  []render.Drawable
	Transform *Transform2D
}

func (m *Manager) NewElement(x, y, w, h, z float32) *Element {
	e := &Element{
		width:    w,
		height:   h,
		children: []render.Drawable{},

		Transform: CreateTransform2D(x, y, z),
	}
	return e
}

func (e *Element) ZIndex() float32 {
	// not sure how this is going to work yet
	// parents must be drawn underneath children (?)
	if e.parent != nil {
		return e.Parent().ZIndex() - e.z
	}
	return e.z
}

// Parent peturns the parent element
func (e *Element) Parent() render.Drawable {
	return e.parent
}

// SetParent sets the parent element
func (e *Element) SetParent(parent render.Drawable) {
	// TODO detach from current parent?
	e.parent = parent
}

// Children returns a list of child elements
func (e *Element) Children() []render.Drawable {
	return e.children
}

func (e *Element) Width() float32 {
	return e.width
}

func (e *Element) Height() float32 {
	return e.height
}

// Attach a child to this element
func (e *Element) Attach(child render.Drawable) {
	e.children = append(e.children, child)
	// set parent?
}

// Detach a child from this element
func (e *Element) Detach(child render.Drawable) {
	// TODO Implement
	//child.Parent = nil
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
	right := e.Transform.Position.X() + e.Width()
	bottom := e.Transform.Position.Y() + e.Height()
	return pos.X() >= e.Transform.Position.X() &&
		pos.Y() >= e.Transform.Position.Y() &&
		pos.X() <= right &&
		pos.Y() <= bottom
}
