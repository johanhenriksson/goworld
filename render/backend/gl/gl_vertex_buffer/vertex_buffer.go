package gl_vertex_buffer

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// glvertexbuf represents an OpenGL Vertex Buffer object
type glvertexbuf struct {
	ID     uint32 /* OpenGL Buffer Identifier */
	Target uint32 /* Target buffer type, defaults to GL_ARRAY_BUFFER */
	Usage  uint32 /* Buffer usage flag, defaults to GL_STATIC_DRAW */
	Size   int    /* Element size in bytes */
}

// New creates a new GL vertex buffer object
func New() vertex.Buffer {
	vbo := &glvertexbuf{
		Target: gl.ARRAY_BUFFER,
		Usage:  gl.STATIC_DRAW,
	}
	gl.GenBuffers(1, &vbo.ID)
	return vbo
}

// NewIndex creates a new GL index buffer object
func NewIndex() *glvertexbuf {
	vbo := &glvertexbuf{
		Target: gl.ELEMENT_ARRAY_BUFFER,
		Usage:  gl.STATIC_DRAW,
	}
	gl.GenBuffers(1, &vbo.ID)
	return vbo
}

// Bind the vertex buffer object
func (vbo *glvertexbuf) Bind() {
	if vbo.ID == 0 {
		panic(fmt.Errorf("cant bind vbo id 0"))
	}
	gl.BindBuffer(vbo.Target, vbo.ID)
}

// Unbind the vertex buffer object
func (vbo *glvertexbuf) Unbind() {
	gl.BindBuffer(vbo.Target, 0)
}

// Delete frees the GPU memory allocated by this vertex buffer. Resets ID and Size to 0
func (vbo *glvertexbuf) Delete() {
	if vbo.ID != 0 {
		gl.DeleteBuffers(1, &vbo.ID)
		vbo.ID = 0
		vbo.Size = 0
	}
}

func (vbo *glvertexbuf) BufferFrom(ptr unsafe.Pointer, size int) {
	vbo.Size = size

	if size <= 0 {
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
func (vbo *glvertexbuf) Buffer(data interface{}) {
	// make sure we've been passed a slice
	t := reflect.TypeOf(data)
	if t.Kind() != reflect.Slice {
		panic(fmt.Errorf("buffered data must be a slice"))
	}

	v := reflect.ValueOf(data)

	// the length of the slice is the number of buffer elements
	elements := v.Len()

	// get byte size of each element, e.g. sizeof(element)
	size := int(t.Elem().Size())

	// get a pointer to the beginning of the array
	ptr := unsafe.Pointer(v.Pointer())

	vbo.BufferFrom(ptr, elements*size)
}
