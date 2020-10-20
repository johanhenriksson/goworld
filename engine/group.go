package engine

import (
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Group can be used to position multiple components relative to a single transform
type Group struct {
	*Transform
	Components []Component
}

// NewGroup creates a new object group
func NewGroup(position, rotation vec3.T) *Group {
	return &Group{
		Transform:  NewTransform(position, rotation, vec3.One),
		Components: []Component{},
	}
}

// Attach a component to the group
func (o *Group) Attach(components ...Component) {
	o.Components = append(o.Components, components...)
}

// Update all components in the group
func (o *Group) Update(dt float32) {
	o.Transform.Update(dt)
	Update(dt, o.Components...)
}

func (o *Group) Collect(pass DrawPass, args DrawArgs) {
	Collect(pass, args.Apply(o.Transform), o.Components...)
}
