package render

import (
    "fmt"
    "github.com/go-gl/gl/v4.1-core/gl"
)

/* Interface for data types that can be uploaded into a vertex buffer object */
type VertexData interface {
    Elements()  int     /* Number of items, usually len(slice) */
    Size()      int     /* Size of each individual element */
}

/* Represents an OpenGL Vertex Buffer object */
type VertexBuffer struct {
    Id          uint32  /* OpenGL Buffer Identifier */
    Elements    int     /* Number of elements/primitives currently stored in the VBO */
    Size        int     /* Element size */

    usage       uint32  /* Buffer usage flag, defaults to GL_STATIC_DRAW */
}

func CreateVertexBuffer() *VertexBuffer {
    vbo := &VertexBuffer {
        usage: gl.STATIC_DRAW,
    }
    gl.GenBuffers(1, &vbo.Id)
    return vbo
}

func (vbo *VertexBuffer) Bind() {
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo.Id)
}

/* Frees the GPU memory allocated by this vertex buffer. Resets Id, Size and Elements to 0 */
func (vbo *VertexBuffer) Delete() {
    gl.DeleteBuffers(1, &vbo.Id)
    *vbo = VertexBuffer { }
}

/* Binds the VBO and buffers data to the GPU */
func (vbo *VertexBuffer) Buffer(vertices VertexData) {
    vbo.Bind()
    vbo.Elements = vertices.Elements()
    vbo.Size     = vertices.Size()

    size := vbo.Size * vbo.Elements
    ptr  := gl.Ptr(vertices)

    // upload got GPU
    gl.BufferData(gl.ARRAY_BUFFER, size, ptr, vbo.usage)

    // debug logging
    fmt.Println("Buffering", vbo.Elements, "elements to buffer", vbo.Id)
    fmt.Println("Element size:", vbo.Size, "bytes")
    fmt.Println("Total size:", size, "bytes")
}

type FloatBuffer []float32

func (vtx FloatBuffer) Elements() int {
    return len(vtx)
}

func (vtx FloatBuffer) Size() int {
    return 4
}
