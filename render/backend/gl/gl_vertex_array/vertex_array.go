package gl_vertex_array

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/backend/gl"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_buffer"
	"github.com/johanhenriksson/goworld/render/vertex"

	ogl "github.com/go-gl/gl/v4.1-core/gl"
)

var activeVAO = 0

const Vertex = "vertex"
const Index = "index"

type BufferMap map[string]vertex.Buffer

// glvertexarray represents an OpenGL Vertex Array Object (VAO)
type glvertexarray struct {
	id       int              // opengl handle
	mode     vertex.Primitive // primitive type
	elements int              // number of elements
	vbos     BufferMap
	index    gl.Type
}

// New creates a new vertex array object. Default primitive is GL_TRIANGLES
func New(primitive vertex.Primitive) vertex.Array {
	vao := &glvertexarray{
		mode:  primitive,
		vbos:  BufferMap{},
		index: gl.None,
	}

	// create vao
	id := uint32(vao.id)
	ogl.GenVertexArrays(1, &id)
	vao.id = int(id)

	// leave it bound
	vao.Bind()
	return vao
}

func (vao *glvertexarray) Indexed() bool {
	return vao.index != gl.None
}

// Delete frees the memory associated with this vertex array object
func (vao *glvertexarray) Delete() {
	if vao.id == 0 {
		return
	}

	// delete vbos
	for _, vbo := range vao.vbos {
		vbo.Delete()
	}

	id := uint32(vao.id)
	ogl.DeleteVertexArrays(1, &id)
	*vao = glvertexarray{}
}

// Bind this vertex array object
func (vao glvertexarray) Bind() error {
	if vao.id == 0 {
		return fmt.Errorf("attempt to bind Vertex Array with ID 0")
	}
	if activeVAO != vao.id {
		ogl.BindVertexArray(uint32(vao.id))
		activeVAO = vao.id
	}
	return nil
}

// Unbind the vertex array
func (vao glvertexarray) Unbind() {
	ogl.BindVertexArray(0)
	activeVAO = 0
}

// Draw all elements in the vertex array
func (vao glvertexarray) Draw() error {
	if vao.elements == 0 {
		// fmt.Println("warning: attempt to draw VAO with length 0")
		return nil
	}

	// draw call
	if err := vao.Bind(); err != nil {
		return err
	}

	if !vao.Indexed() {
		ogl.DrawArrays(uint32(vao.mode), 0, int32(vao.elements))
	} else {
		ogl.DrawElements(uint32(vao.mode), int32(vao.elements), uint32(vao.index), nil)
	}

	return nil
}

// Buffer vertex data to the GPU
func (vao *glvertexarray) Buffer(name string, data any) {
	vao.Bind()

	vbo, exists := vao.vbos[name]
	if !exists {
		// create new buffer
		if name == Index {
			vbo = gl_vertex_buffer.NewIndex()
		} else if name == Vertex {
			vbo = gl_vertex_buffer.New()
		} else {
			panic(fmt.Sprintf("illegal buffer name: %s", name))
		}
		vao.vbos[name] = vbo
	}

	vbo.Buffer(data)
}

func (vao *glvertexarray) SetPointers(pointers vertex.Pointers) {
	// we need to have a vertex buffer before binding any pointers.
	// so if it does not exist yet, we create it here. this helps avoid bugs
	vao.Bind()
	vbo, exists := vao.vbos[Vertex]
	if !exists {
		vbo = gl_vertex_buffer.New()
		vao.vbos[Vertex] = vbo
	}
	vbo.Bind()
	gl.EnablePointers(pointers)
}

func (vao *glvertexarray) SetIndexSize(size int) {
	switch size {
	case 0:
		vao.index = gl.None
	case 1:
		vao.index = gl.UInt8
	case 2:
		vao.index = gl.UInt16
	case 4:
		vao.index = gl.UInt32
	default:
		panic(fmt.Sprintf("illegal index size: %d", size))
	}
}

func (vao *glvertexarray) SetElements(count int) {
	vao.elements = count
}
