package engine

import (
    "github.com/johanhenriksson/goworld/render"
)

/** Game object */
type Object struct {
    *Transform
    Components  []Component
    Children    []*Object
}

func NewObject(x,y,z float32) *Object {
    return &Object {
        Transform: CreateTransform(x,y,z),
        Components: []Component { },
        Children: []*Object { },
    }
}

func (o *Object) Attach(component Component) {
    o.Components = append(o.Components, component)
}

func (o *Object) Draw(args render.DrawArgs) {
    /* Apply transform */
    args.Transform = o.Transform.Matrix.Mul4(args.Transform)

    /* Draw components */
    args.Shader.Matrix4f("model", &args.Transform[0])
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