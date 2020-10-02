package engine

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Object is the basic component of a scene. It has a transform, a list of components and optionally child objects.
type Object struct {
	*Transform
	Components []Component
	Children   []*Object
}

// NewObject creates a new object in the scene.
func NewObject(position vec3.T) *Object {
	return &Object{
		Transform:  CreateTransform(position),
		Components: []Component{},
		Children:   []*Object{},
	}
}

// Attach a component to the object
func (o *Object) Attach(component Component) {
	o.Components = append(o.Components, component)
}

// Draw the object, its components and its children.
func (o *Object) Draw(args render.DrawArgs) {
	// apply transforms
	args.Transform = o.Transform.Matrix.Mul(&args.Transform)
	args.MVP = args.VP.Mul(&args.Transform)

	// draw components
	for _, comp := range o.Components {
		comp.Draw(args)
	}

	// draw children
	for _, child := range o.Children {
		child.Draw(args)
	}
}

// Update the object, its components and its children.
func (o *Object) Update(dt float32) {
	o.Transform.Update(dt)

	// update components
	for _, comp := range o.Components {
		comp.Update(dt)
	}

	// update children
	for _, child := range o.Children {
		child.Update(dt)
	}
}
