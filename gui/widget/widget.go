package widget

import (
	"fmt"

	"github.com/johanhenriksson/goworld/gui/dimension"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render"
)

type T interface {
	Key() string

	// Properties returns a pointer to the components property struct.
	// The pointer is used to compare the states when deciding if the component needs to be updated.
	Props() any

	// Update replaces the components property struct.
	Update(any)

	// Destroy releases any resources associated with the component.
	// Attempting to draw a destroyed component will cause a panic.
	Destroy()

	// Destroyed indicates whether the component has been destroyed or not.
	Destroyed() bool

	// Size returns the actual size of the element in pixels
	Size() vec2.T

	// Position returns the current position of the element relative to its parent
	Position() vec2.T

	SetSize(vec2.T)
	SetPosition(vec2.T)

	Width() dimension.T
	Height() dimension.T

	Children() []T
	SetChildren([]T)

	Arrange(vec2.T) vec2.T

	// Draw the widget. This should only be called by the GUI Draw Pass
	// Calling Draw() will instantiate any required GPU resources prior to drawing.
	// Attempting to draw a destroyed component will cause a panic.
	Draw(render.Args)
}

type widget struct {
	id        int
	key       string
	size      vec2.T
	position  vec2.T
	destroyed bool
}

func New(key string) T {
	return &widget{
		key: key,
	}
}

func (w *widget) Key() string          { return w.key }
func (w *widget) Position() vec2.T     { return w.position }
func (w *widget) Size() vec2.T         { return w.size }
func (w *widget) Width() dimension.T   { return dimension.Auto() }
func (w *widget) Height() dimension.T  { return dimension.Auto() }
func (w *widget) Destroyed() bool      { return w.destroyed }
func (w *widget) SetPosition(p vec2.T) { w.position = p }
func (w *widget) SetSize(s vec2.T)     { w.size = s }
func (w *widget) Destroy()             { w.destroyed = true }

func (w *widget) Props() any {
	panic("widget.Props() must be implemented")
}

func (w *widget) Update(any) {
	panic("widget.Update() must be implemented")
}

func (w *widget) Children() []T     { return nil }
func (w *widget) SetChildren(c []T) {}

func (w *widget) Draw(render.Args) {
	// base widget Draw() should be called ahead of overridden draws

	if w.Destroyed() {
		panic(fmt.Sprintf("attempt to draw destroyed widget %s", w.key))
	}
}

func (w *widget) Arrange(space vec2.T) vec2.T {
	w.SetSize(space)
	return space
}

func (w *widget) Measure(space vec2.T) vec2.T {
	// this could also take in to account whether the element wants to grow or shrink etc
	// also dimensions of the children
	return vec2.New(
		w.Width().Resolve(space.X),
		w.Height().Resolve(space.Y))
}
