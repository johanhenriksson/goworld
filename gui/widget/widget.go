package widget

import (
	"github.com/johanhenriksson/goworld/math/vec2"

	"github.com/kjk/flex"
)

type T interface {
	Key() string

	// Properties returns a pointer to the components property struct.
	// The pointer is used to compare the states when deciding if the component needs to be updated.
	Props() any

	// Update replaces the components property struct.
	Update(any)

	// Size returns the actual size of the element in pixels
	Size() vec2.T

	// Position returns the current position of the element relative to its parent
	Position() vec2.T

	Children() []T
	SetChildren([]T)

	Flex() *flex.Node

	// Draw the widget. This should only be called by the GUI Draw Pass
	// Calling Draw() will instantiate any required GPU resources prior to drawing.
	// Attempting to draw a destroyed component will cause a panic.
	Draw(DrawArgs, *QuadBuffer)
}
