package render

import (
	"fmt"
	"reflect"
	"unsafe"

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
func (vbo *VertexBuffer) Bind() {
	if vbo.ID == 0 {
		panic(fmt.Errorf("cant bind vbo id 0"))
	}
	gl.BindBuffer(vbo.Target, vbo.ID)
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

type BufferCommand struct {
	Elements int
	Size     int
	Source   unsafe.Pointer
}

func (vbo *VertexBuffer) BufferFrom(elements, size int, ptr unsafe.Pointer) {
	vbo.Elements = elements
	vbo.Size = size * elements

	if elements <= 0 {
		return
	}

	// buffer data to GPU
	vbo.Bind()
	gl.BufferData(vbo.Target, vbo.Size, ptr, vbo.Usage)

	// check actual size in GPU memory
	gpuSize := int32(0)
	gl.GetBufferParameteriv(vbo.Target, gl.BUFFER_SIZE, &gpuSize)
	if int(gpuSize) != vbo.Size {
		panic(fmt.Errorf("failed to buffer data to vbo #%d, expected size %d bytes, actual: %d bytes",
			vbo.ID, vbo.Size, gpuSize))
	}
}

// Buffer data to GPU memory
func (vbo *VertexBuffer) Buffer(data interface{}) int {
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Slice {
		panic(fmt.Errorf("buffered data must be a slice"))
	}
	v := reflect.ValueOf(data)
	elements := v.Len()
	size := int(t.Elem().Size())
	ptr := unsafe.Pointer(v.Pointer())

	vbo.BufferFrom(elements, size, ptr)
	return elements
}
