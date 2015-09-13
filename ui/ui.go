package ui;

import (
    "unsafe"
    "github.com/go-gl/gl/v4.1-core/gl"
    mgl "github.com/go-gl/mathgl/mgl32"
)

type Drawable interface {
    Draw(DrawArgs)
}

type DrawArgs struct {
    Viewport    mgl.Mat4
    Transform   mgl.Mat4
}

/** UI Color type */
type Color struct {
    R, G, B, A  float32
}

/** Convinience method - returns a new color struct */
func NewColor(r,g,b,a float32) Color {
    return Color { R: r, G: g, B: b, A: a, }
}

/** Color vertex. Used in solid-color elements */
type ColorVertex struct {
    X, Y, Z     float32 // 12 bytes
    Color               // 16 bytes
} // 28 bytes

type ColorVertices []ColorVertex
func (buffer ColorVertices) Elements() int { return len(buffer) }
func (buffer ColorVertices) Size()     int { return 28 }
func (buffer ColorVertices) GLPtr()    unsafe.Pointer { return gl.Ptr(buffer) }

/** ImageVertex */
type ImageVertex struct {
    X, Y, Z     float32 // 12 bytes
    Tx, Ty      float32 // 8 bytes
}

type ImageVertices []ImageVertex
func (buffer ImageVertices) Elements() int { return len(buffer) }
func (buffer ImageVertices) Size()     int { return 28 }
func (buffer ImageVertices) GLPtr()    unsafe.Pointer { return gl.Ptr(buffer) }
