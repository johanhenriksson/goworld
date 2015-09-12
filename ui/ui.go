package ui;

import (
    "unsafe"
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/geometry"
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

type Color struct {
    R, G, B, A  float32
}

type ColorVertex struct {
    X, Y, Z     float32 // 12 bytes
    Color
} // 28 bytes

type ColorVertices []ColorVertex

func (buffer ColorVertices) Elements() int { return len(buffer) }
func (buffer ColorVertices) Size()     int { return 28 }
func (buffer ColorVertices) GLPtr()    unsafe.Pointer { return gl.Ptr(buffer) }

type Quad struct {
    TopLeft     ColorVertex
    TopRight    ColorVertex
    BottomLeft  ColorVertex
    BottomRight ColorVertex
    Material    *render.Material
    vao         *geometry.VertexArray
    vbo         *geometry.VertexBuffer
}

func NewQuad(mat *render.Material, color Color, w,h,z,r,g,b,a float32) *Quad {
    q := &Quad {
        Material: mat,
        TopLeft: ColorVertex {
            X: 0,
            Y: h,
            Z: z,
            Color: color,
        },
        TopRight: ColorVertex {
            X: w,
            Y: h,
            Z: z,
            Color: color,
        },
        BottomLeft: ColorVertex {
            X: 0,
            Y: 0,
            Z: z,
            Color: color,
        },
        BottomRight: ColorVertex {
            X: w,
            Y: 0,
            Z: z,
            Color: color,
        },
        vao: geometry.CreateVertexArray(),
        vbo: geometry.CreateVertexBuffer(),
    }
    q.compute()
    return q
}

func (q *Quad) compute() {
    vtx := ColorVertices {
        q.BottomLeft, q.TopRight, q.TopLeft,
        q.BottomLeft, q.BottomRight, q.TopRight,
    }
    q.vao.Length = 6
    q.vao.Bind()
    q.vbo.Buffer(vtx)
    q.Material.Setup()
}

func (q *Quad) Draw(args DrawArgs) {
    q.Material.Use()
    q.Material.Shader.Matrix4f("model", &args.Transform[0])
    q.Material.Shader.Matrix4f("viewport", &args.Viewport[0])
    q.vao.Draw()
}
