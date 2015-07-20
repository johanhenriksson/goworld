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

type VertexData interface {
    Elements() int
    Size() int
    GLPtr() unsafe.Pointer
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
    fmt.Println("Buffering", vbo.Elements, "elements to buffer", vbo.Id)
    fmt.Println("Element size:", vbo.Size, "bytes.")
    fmt.Println("Total size:", size, "bytes.")
	gl.BufferData(gl.ARRAY_BUFFER, size, vertices.GLPtr(), gl.STATIC_DRAW)
}
