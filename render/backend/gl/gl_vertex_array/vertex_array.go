package gl_vertex_array

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_buffer"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/vertex_array"
	"github.com/johanhenriksson/goworld/render/vertex_buffer"

	ogl "github.com/go-gl/gl/v4.1-core/gl"
)

type BufferMap map[string]vertex_buffer.T

// glvertexarray represents an OpenGL Vertex Array Object (VAO)
type glvertexarray struct {
	ID     int              /* OpenGL Vertex Array identifier */
	Type   render.Primitive /* Primitive type */
	Length int              /* Number of verticies */

	vbos  BufferMap
	index gl.Type
}

// New creates a new vertex array object. Default primitive is GL_TRIANGLES
func New(primitive render.Primitive) vertex_array.T {
	vao := &glvertexarray{
		Type: primitive,
		vbos: BufferMap{},
	}

	// create vao
	id := uint32(vao.ID)
	ogl.GenVertexArrays(1, &id)
	vao.ID = int(id)

	// leave it bound
	vao.Bind()
	return vao
}

func (vao *glvertexarray) Indexed() bool {
	return vao.index != ogl.NONE
}

// Delete frees the memory associated with this vertex array object
func (vao *glvertexarray) Delete() {
	if vao.ID == 0 {
		return
	}

	// delete vbos
	for _, vbo := range vao.vbos {
		vbo.Delete()
	}

	id := uint32(vao.ID)
	ogl.DeleteVertexArrays(1, &id)
	*vao = glvertexarray{}
}

// Bind this vertex array object
func (vao glvertexarray) Bind() {
	if vao.ID == 0 {
		fmt.Println("warning: attempt to bind Vertex Array with ID 0")
	}
	ogl.BindVertexArray(uint32(vao.ID))
}

// Unbind the vertex array
func (vao glvertexarray) Unbind() {
	ogl.BindVertexArray(0)
}

// Draw all elements in the vertex array
func (vao glvertexarray) Draw() {
	if vao.Length == 0 {
		// fmt.Println("warning: attempt to draw VAO with length 0")
		return
	}

	// draw call
	vao.Bind()

	if !vao.Indexed() {
		ogl.DrawArrays(uint32(vao.Type), 0, int32(vao.Length))
	} else {
		ogl.DrawElements(uint32(vao.Type), int32(vao.Length), uint32(vao.index), nil)
	}
}

func (vao *glvertexarray) SetIndexType(t types.Type) {
	// todo: get rid of this later
	vao.index = gl.TypeCast(t)
}

// Buffer vertex data to the GPU
func (vao *glvertexarray) Buffer(name string, data interface{}) {
	if name == "index" {
		// todo: set index type
		// then get rid of SetIndexType
	}

	vao.Bind()

	vbo, exists := vao.vbos[name]
	if !exists {
		// create new buffer
		vbo = gl_vertex_buffer.New()
		vao.vbos[name] = vbo
	}

	// buffer data to vbo
	elements := vbo.Buffer(data)

	// update number of elements
	if !vao.Indexed() || name == "index" {
		vao.Length = elements
	}
}

func (vao *glvertexarray) BufferTo(pointers shader.Pointers, data interface{}) {
	name := pointers.BufferString()

	vao.Bind()

	vbo, exists := vao.vbos[name]
	if !exists {
		// create new buffer
		vbo = gl_vertex_buffer.New()
		vao.vbos[name] = vbo
	}

	// buffer data to vbo
	elements := vbo.Buffer(data)

	// update number of elements
	if !vao.Indexed() {
		vao.Length = elements
	}

	pointers.Enable()
}
