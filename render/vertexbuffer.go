package render

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// VertexBuffer represents an OpenGL Vertex Buffer object
type VertexBuffer struct {
	ID       uint32 /* OpenGL Buffer Identifier */
	Target   uint32 /* Target buffer type, defaults to GL_ARRAY_BUFFER */
	Usage    uint32 /* Buffer usage flag, defaults to GL_STATIC_DRAW */
	Elements int    /* Number of verticies/elements currently stored in the VBO */
	Size     int    /* Element size in bytes */
}

// CreateVertexBuffer creates a new GL vertex buffer object
func CreateVertexBuffer() *VertexBuffer {
	vbo := &VertexBuffer{
		Target: gl.ARRAY_BUFFER,
		Usage:  gl.STATIC_DRAW,
	}
	gl.GenBuffers(1, &vbo.ID)
	return vbo
}

// CreateIndexBuffer creates a new GL index buffer object
func CreateIndexBuffer() *VertexBuffer {
	vbo := &VertexBuffer{
		Target: gl.ELEMENT_ARRAY_BUFFER,
		Usage:  gl.STATIC_DRAW,
	}
	gl.GenBuffers(1, &vbo.ID)
	return vbo
}

// Bind the vertex buffer object
func (vbo *VertexBuffer) Bind() error {
	if vbo.ID == 0 {
		return fmt.Errorf("Cannot bind buffer id 0")
	}
	gl.BindBuffer(vbo.Target, vbo.ID)
	return nil
}

// Unbind the vertex buffer object
func (vbo *VertexBuffer) Unbind() {
	gl.BindBuffer(vbo.Target, 0)
}

// Delete frees the GPU memory allocated by this vertex buffer. Resets ID and Size to 0
func (vbo *VertexBuffer) Delete() {
	if vbo.ID != 0 {
		gl.DeleteBuffers(1, &vbo.ID)
		vbo.ID = 0
		vbo.Size = 0
	}
}

// Buffer data to GPU memory
func (vbo *VertexBuffer) Buffer(vertices VertexData) error {
	// bind buffer
	err := vbo.Bind()
	if err != nil {
		return err
	}

	// buffer data to GPU
	size := vertices.Size() * vertices.Elements()
	ptr := gl.Ptr(vertices)
	gl.BufferData(vbo.Target, size, ptr, vbo.Usage)

	// check actual size in GPU memory
	gpuSize := int32(0)
	gl.GetBufferParameteriv(vbo.Target, gl.BUFFER_SIZE, &gpuSize)
	if int(gpuSize) != size {
		return fmt.Errorf("Failed buffering data to buffer #%d, expected size %d bytes, actual: %d bytes",
			vbo.ID, size, gpuSize)
	}

	vbo.Size = vertices.Size()
	vbo.Elements = vertices.Elements()

	// debug logging
	fmt.Printf("[VBO %d] Buffered %d x %d = %d bytes\n", vbo.ID, vbo.Size, vbo.Elements, size)
	return nil
}
