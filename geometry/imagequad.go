package geometry

import (
    "github.com/johanhenriksson/goworld/render"
)

/** Not exactly a quad anymore is it? */
type ImageQuad struct {
    Material    *render.Material
    TopLeft     ImageVertex
    TopRight    ImageVertex
    BottomLeft  ImageVertex
    BottomRight ImageVertex
    vao         *VertexArray
    vbo         *VertexBuffer
}

func NewImageQuad(mat *render.Material, w,h,z float32) *ImageQuad {
    q := &ImageQuad {
        Material:    mat,
        TopLeft:     ImageVertex { X: 0, Y: h, Z: z, Tx: 0, Ty: 0, },
        TopRight:    ImageVertex { X: w, Y: h, Z: z, Tx: 1, Ty: 0, },
        BottomLeft:  ImageVertex { X: 0, Y: 0, Z: z, Tx: 0, Ty: 1, },
        BottomRight: ImageVertex { X: w, Y: 0, Z: z, Tx: 1, Ty: 1, },
        vao: CreateVertexArray(),
        vbo: CreateVertexBuffer(),
    }
    q.compute()
    return q
}


func (q *ImageQuad) compute() {
    vtx := ImageVertices {
        q.BottomLeft, q.TopRight, q.TopLeft,
        q.BottomLeft, q.BottomRight, q.TopRight,
    }

    /* Setup VAO */
    q.vao.Length = int32(len(vtx))
    q.vao.Bind()
    q.vbo.Buffer(vtx)
    q.Material.Setup()
}

func (q *ImageQuad) Draw(args render.DrawArgs) {
    q.Material.Use()
    q.Material.Shader.Matrix4f("model", &args.Transform[0])
    q.Material.Shader.Matrix4f("viewport", &args.Viewport[0])
    q.vao.Draw()
}
