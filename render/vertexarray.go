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
func CreateVertexArray(primitive GLPrimitive, buffers ...string) *VertexArray {
	vao := &VertexArray{
		Type: primitive,
		vbos: BufferMap{},
	}

	// create vao
	id := uint32(vao.ID)
	gl.GenVertexArrays(1, &id)
	vao.ID = int(id)

	// create buffers
	vao.Bind()
	for _, buffer := range buffers {
		vao.AddBuffer(buffer)
	}
	return vao
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
		fmt.Println("Warning: Attempt to bind Vertex Array with ID 0")
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
		fmt.Println("Warning: Attempt to draw VAO with length 0")
	}

	// draw call
	vao.Bind()

	if vao.index == gl.NONE {
		gl.DrawArrays(uint32(vao.Type), 0, int32(vao.Length))
	} else {
		gl.DrawElements(uint32(vao.Type), int32(vao.Length), uint32(vao.index), nil)
	}
}

// AddBuffer adds a named buffer to the VAO.
func (vao *VertexArray) AddBuffer(name string) *VertexBuffer {
	if name == "index" {
		panic("index is reserved for index buffers")
	}
	if vbo, exists := vao.vbos[name]; exists {
		return vbo
	}

	// set up vertex array pointers for this buffer
	vao.Bind()

	// create new vbo
	vbo := CreateVertexBuffer()
	vbo.Bind()

	// store reference & return vbo object
	vao.vbos[name] = vbo
	return vbo
}

// AddIndexBuffer adds an index buffer to the VAO.
func (vao *VertexArray) AddIndexBuffer(datatype GLType) *VertexBuffer {
	// set up vertex array pointers for this buffer
	vao.Bind()
	vao.index = datatype

	// create new vbo
	vbo := CreateIndexBuffer()
	vbo.Bind()

	// store reference & return vbo object
	vao.vbos["index"] = vbo
	return vbo
}

// Buffer vertex data to the GPU
func (vao *VertexArray) Buffer(name string, data VertexData) error {
	vbo, exists := vao.vbos[name]
	if !exists {
		panic(fmt.Sprintf("Unknown VBO: %s", name))
	}

	if data.Elements() == 0 {
		vao.Length = 0
		return nil
	}

	vao.Bind()
	vao.Length = data.Elements()
	return vbo.Buffer(data)
}
