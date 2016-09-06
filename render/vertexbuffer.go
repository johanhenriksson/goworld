package render

import (
    "fmt"
    "github.com/go-gl/gl/v4.1-core/gl"
)

type VertexBuffer struct {
    Id          uint32
    Elements    int
    Size        int
}

type VertexData interface {
    Elements() int
    Size() int
}

func CreateVertexBuffer() *VertexBuffer {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
    return &VertexBuffer {
        Id: vbo,
    }
}

func (vbo *VertexBuffer) Bind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo.Id)
}


func (vbo *VertexBuffer) Buffer(vertices VertexData) {
    vbo.Bind()
    vbo.Elements = vertices.Elements()
    vbo.Size     = vertices.Size()
    size := vbo.Size * vbo.Elements
    ptr  := gl.Ptr(vertices)
    fmt.Println("Buffering", vbo.Elements, "elements to buffer", vbo.Id)
    fmt.Println("Element size:", vbo.Size, "bytes.")
    fmt.Println("Total size:", size, "bytes.")
	gl.BufferData(gl.ARRAY_BUFFER, size, ptr, gl.STATIC_DRAW)
}

type FloatBuffer []float32

func (vtx FloatBuffer) Elements() int {
    return len(vtx)
}

func (vtx FloatBuffer) Size() int {
    return 4
}
