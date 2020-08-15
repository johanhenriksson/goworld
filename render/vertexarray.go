package render

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// VertexArray represents an OpenGL Vertex Array Object (VAO)
type VertexArray struct {
	ID     uint32 /* OpenGL Vertex Array identifier */
	Type   uint32 /* Primitive type */
	Length int32  /* Number of verticies */
}

// CreateVertexArray creates a new vertex array object. Default primitive is GL_TRIANGLES
func CreateVertexArray() *VertexArray {
	vao := &VertexArray{
		Type: gl.TRIANGLES,
	}
	gl.GenVertexArrays(1, &vao.ID)
	return vao
}

// Delete frees the memory associated with this vertex array object
func (vao *VertexArray) Delete() {
	if vao.ID != 0 {
		gl.DeleteVertexArrays(1, &vao.ID)
		*vao = VertexArray{}
	}
}

// Bind this vertex array object
func (vao VertexArray) Bind() error {
	if vao.ID == 0 {
		return fmt.Errorf("Cannot bind Vertex Array id 0")
	}
	gl.BindVertexArray(vao.ID)
	return nil
}

// Unbind the vertex array
func (vao VertexArray) Unbind() {
	gl.BindVertexArray(0)
}

// DrawElements draws a range of elements in the vertex array
func (vao VertexArray) DrawElements() error {
	// bind vertex array
	err := vao.Bind()
	if err != nil {
		return err
	}

	// draw call
	gl.DrawArrays(vao.Type, 0, vao.Length)
	return nil
}

// DrawIndexed draws using the index buffer
func (vao VertexArray) DrawIndexed() error {
	// bind vertex array
	err := vao.Bind()
	if err != nil {
		return err
	}

	gl.DrawElements(vao.Type, vao.Length, gl.UNSIGNED_INT, nil)
	return nil
}
