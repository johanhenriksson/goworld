package render

import (
    "fmt"
    "github.com/go-gl/gl/v4.1-core/gl"
)

/* Represents an OpenGL Vertex Buffer object */
type VertexBuffer struct {
    Id          uint32  /* OpenGL Buffer Identifier */
    Target      uint32  /* Target buffer type, defaults to GL_ARRAY_BUFFER */
    Usage       uint32  /* Buffer usage flag, defaults to GL_STATIC_DRAW */
    Elements    int     /* Number of verticies/elements currently stored in the VBO */
    Size        int     /* Element size in bytes */
}

/* Create a new Vertex buffer object and allocate a matching OpenGL buffer */
func CreateVertexBuffer() *VertexBuffer {
    vbo := &VertexBuffer {
        Target: gl.ARRAY_BUFFER,
        Usage:  gl.STATIC_DRAW,
    }
    gl.GenBuffers(1, &vbo.Id)
    return vbo
}

/* Binds the vertex buffer object */
func (vbo *VertexBuffer) Bind() error {
    if vbo.Id == 0 {
        return fmt.Errorf("Cannot bind buffer id 0")
    }
    gl.BindBuffer(vbo.Target, vbo.Id)
    return nil
}

/* Frees the GPU memory allocated by this vertex buffer. Resets Id, Size and Elements to 0 */
func (vbo *VertexBuffer) Delete() {
    if vbo.Id != 0 {
        gl.DeleteBuffers(1, &vbo.Id)
        *vbo = VertexBuffer { }
    }
}

/* Binds the VBO and buffers data to the GPU */
func (vbo *VertexBuffer) Buffer(vertices VertexData) error {
    // bind buffer
    err := vbo.Bind()
    if err != nil {
        return err
    }

    // buffer data to GPU
    size := vertices.Size() * vertices.Elements()
    ptr  := gl.Ptr(vertices)
    gl.BufferData(vbo.Target, size, ptr, vbo.Usage)

    // check actual size in GPU memory
    var gpuSize int32 = 0
    gl.GetBufferParameteriv(vbo.Target, gl.BUFFER_SIZE, &gpuSize)
    if int(gpuSize) != size {
        return fmt.Errorf("Failed buffering data to buffer #%d, expected size %d bytes, actual: %d bytes",
            vbo.Id, size, gpuSize)
    }

    vbo.Size     = vertices.Size()
    vbo.Elements = vertices.Elements()

    // debug logging
    fmt.Printf("[VBO %d] Buffered %d x %d = %d bytes\n", vbo.Id, vbo.Size, vbo.Elements, size)

    return nil
}
