package engine

import (
	"reflect"

	"github.com/johanhenriksson/goworld/render"
)

/** Game object */
type Object struct {
	*Transform
	Scene      *Scene
	Components []Component
	Children   []*Object
}

func (s *Scene) NewObject(x, y, z float32) *Object {
	return &Object{
		Transform:  CreateTransform(x, y, z),
		Scene:      s,
		Components: []Component{},
		Children:   []*Object{},
	}
}

func (o *Object) Attach(component Component) {
	o.Components = append(o.Components, component)
}

func (o *Object) Draw(args render.DrawArgs) {
	/* Apply transform */
	args.Transform = o.Transform.Matrix.Mul4(args.Transform)

	/* Draw components */
	args.MVP = args.VP.Mul4(args.Transform)
	args.Shader.Mat4f("mvp", args.MVP)

	// model matrix is required to calculate vertex normals during the geometry pass
	if args.Pass == "geometry" {
		args.Shader.Mat4f("model", args.Transform)
		args.Shader.Mat4f("view", args.View)
		args.Shader.Mat4f("projection", args.Projection)
	}

	for _, comp := range o.Components {
		comp.Draw(args)
	}

	/* Draw children */
	for _, child := range o.Children {
		child.Draw(args)
	}
}

func (o *Object) Update(dt float32) {
	o.Transform.Update(dt)

	/* Update components */
	for _, comp := range o.Components {
		comp.Update(dt)
	}

	/* Update children */
	for _, child := range o.Children {
		child.Update(dt)
	}
}

func (o *Object) GetComponent(component Component) (Component, bool) {
	t := reflect.TypeOf(component)
	for _, c := range o.Components {
		if c.Type() == t {
			return c, true
		}
	}
	return component, false
}
