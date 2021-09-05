package render

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type BufferMap map[string]*VertexBuffer

// VertexArray represents an OpenGL Vertex Array Object (VAO)
type VertexArray struct {
	ID     int         /* OpenGL Vertex Array identifier */
	Type   GLPrimitive /* Primitive type */
	Length int         /* Number of verticies */

	vbos  BufferMap
	index GLType
}

// CreateVertexArray creates a new vertex array object. Default primitive is GL_TRIANGLES
func CreateVertexArray(primitive GLPrimitive) *VertexArray {
	vao := &VertexArray{
		Type: primitive,
		vbos: BufferMap{},
	}

	// create vao
	id := uint32(vao.ID)
	gl.GenVertexArrays(1, &id)
	vao.ID = int(id)

	// leave it bound
	vao.Bind()
	return vao
}

func (vao *VertexArray) Indexed() bool {
	return vao.index != gl.NONE
}

// Delete frees the memory associated with this vertex array object
func (vao *VertexArray) Delete() {
	if vao.ID == 0 {
		return
	}

	// delete vbos
	for _, vbo := range vao.vbos {
		vbo.Delete()
	}

	id := uint32(vao.ID)
	gl.DeleteVertexArrays(1, &id)
	*vao = VertexArray{}
}

// Bind this vertex array object
func (vao VertexArray) Bind() {
	if vao.ID == 0 {
		fmt.Println("warning: attempt to bind Vertex Array with ID 0")
	}
	gl.BindVertexArray(uint32(vao.ID))
}

// Unbind the vertex array
func (vao VertexArray) Unbind() {
	gl.BindVertexArray(0)
}

// Draw all elements in the vertex array
func (vao VertexArray) Draw() {
	if vao.Length == 0 {
		// fmt.Println("warning: attempt to draw VAO with length 0")
		return
	}

	// draw call
	vao.Bind()

	if !vao.Indexed() {
		gl.DrawArrays(uint32(vao.Type), 0, int32(vao.Length))
	} else {
		gl.DrawElements(uint32(vao.Type), int32(vao.Length), uint32(vao.index), nil)
	}
}

func (vao *VertexArray) SetIndexType(t GLType) {
	// get rid of this later
	vao.index = t
}

// Buffer vertex data to the GPU
func (vao *VertexArray) Buffer(name string, data interface{}) {
	if name == "index" {
		// todo: set index type
		// then get rid of SetIndexType
	}

	vao.Bind()

	vbo, exists := vao.vbos[name]
	if !exists {
		// create new buffer
		vbo = CreateVertexBuffer()
		vao.vbos[name] = vbo
	}

	// buffer data to vbo
	elements := vbo.Buffer(data)

	// update number of elements
	if !vao.Indexed() || name == "index" {
		vao.Length = elements
	}
}

func (vao *VertexArray) BufferTo(pointers Pointers, data interface{}) {
	name := pointers.BufferString()

	vao.Bind()

	vbo, exists := vao.vbos[name]
	if !exists {
		// create new buffer
		vbo = CreateVertexBuffer()
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
