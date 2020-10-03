package engine

import (
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Object is the basic component of a scene. It has a transform, a list of components and optionally child objects.
type Object struct {
	*Transform
	Components []Component
}

// NewObject creates a new object in the scene.
func NewObject(position, rotation vec3.T) *Object {
	return &Object{
		Transform:  CreateTransform(position, rotation, vec3.One),
		Components: []Component{},
	}
}

// Attach a component to the object
func (o *Object) Attach(components ...Component) {
	o.Components = append(o.Components, components...)
}

// Draw the object, its components and its children.
func (o *Object) Draw(args DrawArgs) {
	local := args.Apply(o.Transform)
	for _, comp := range o.Components {
		comp.Draw(local)
	}
}

// Update the object, its components and its children.
func (o *Object) Update(dt float32) {
	o.Transform.Update(dt)
	for _, comp := range o.Components {
		comp.Update(dt)
	}
}

func (o *Object) Base() *Object {
	return o
}
