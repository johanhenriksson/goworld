package geometry

import (
    "fmt"
    "unsafe"
    "github.com/go-gl/gl/v4.1-core/gl"
)

type VertexBuffer struct {
    Id          uint32
    Elements    int
    Size        int
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

func (vbo *VertexBuffer) Buffer(vertices []Vertex) {
    vbo.Elements = len(vertices)
    vbo.Size     = int(unsafe.Sizeof(vertices[0]))
    size := vbo.Size * vbo.Elements
    fmt.Println("Buffering", vbo.Elements, ", el size:", vbo.Size, "total:", size)
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(vertices), gl.STATIC_DRAW)
}
