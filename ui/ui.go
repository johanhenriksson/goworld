package ui;

import (
    mgl "github.com/go-gl/mathgl/mgl32"
)

type Drawable interface {
    Draw(DrawArgs)
}

type DrawArgs struct {
    Viewport    mgl.Mat4
    Transform   mgl.Mat4
}
