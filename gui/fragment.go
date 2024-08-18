package gui

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/gui/node"
)

type FragmentPosition int

const FragmentLast FragmentPosition = 0
const FragmentFirst FragmentPosition = 1

// A Fragment contains a UI subtree to be rendered as a child of a given element
type Fragment interface {
	object.Component
	Render() node.T

	Slot() string
	Position() FragmentPosition
}

type FragmentArgs struct {
	Slot     string
	Position FragmentPosition
	Render   node.RenderFunc
}

type fragment struct {
	object.Component
	FragmentArgs
}

func NewFragment(pool object.Pool, args FragmentArgs) Fragment {
	return object.NewComponent(pool, &fragment{
		FragmentArgs: args,
	})
}

func (f *fragment) Name() string { return "UIFragment" }

func (f *fragment) Slot() string {
	return f.FragmentArgs.Slot
}

func (f *fragment) Position() FragmentPosition {
	return f.FragmentArgs.Position
}

func (f *fragment) Render() node.T {
	return f.FragmentArgs.Render()
}
