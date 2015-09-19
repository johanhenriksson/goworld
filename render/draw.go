package render

import (
    mgl "github.com/go-gl/mathgl/mgl32"
)

/** UI Component render interface */
type Drawable interface {
    Draw(DrawArgs)

    ZIndex() float32

    /* Render tree */
    Parent() Drawable
    SetParent(Drawable)
    Children() []Drawable
}

/** Passed to Drawables on render */
type DrawArgs struct {
    Viewport    mgl.Mat4
    Transform   mgl.Mat4
    Shader      *ShaderProgram
}
