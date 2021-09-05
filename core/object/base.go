package object

import (
	"github.com/johanhenriksson/goworld/core/transform"
)

// base contains the minimum required functionality for implementing object components
type base struct {
	parent T
	active bool
}

// NewComponent creates a new component base object
func NewComponent() Component {
	return &base{
		parent: nil,
		active: true,
	}
}

func (b *base) Name() string { return "Component" }

// Object refers to the parent object
func (b *base) Object() T { return b.parent }

// SetObject sets the parent object
func (b *base) SetObject(o T) {
	b.parent = o
}

// Update the component. Called every frame if the component is active.
func (b *base) Update(dt float32) {
	// propagating the update to the linked object causes infinite recursion
	// ... do nothing ...
}

func (b *base) Active() bool {
	return b.active
}

func (b *base) SetActive(active bool) {
	b.active = active
}

func (b *base) Transform() transform.T {
	return b.parent.Transform()
}
