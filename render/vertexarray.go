package render

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

/* Represents an OpenGL Vertex Array Object (VAO) */
type VertexArray struct {
	Id     uint32 /* OpenGL Vertex Array identifier */
	Type   uint32 /* Primitive type */
	Length int32  /* Number of verticies */
}

/* Create a new vertex array object. Default primitive is GL_TRIANGLES */
func CreateVertexArray() *VertexArray {
	vao := &VertexArray{
		Type: gl.TRIANGLES,
	}
	gl.GenVertexArrays(1, &vao.Id)
	return vao
}

/* Frees the memory associated with this vertex array object */
func (vao *VertexArray) Delete() {
	if vao.Id != 0 {
		gl.DeleteVertexArrays(1, &vao.Id)
		*vao = VertexArray{}
	}
}

/* Binds the vertex array */
func (vao VertexArray) Bind() error {
	if vao.Id == 0 {
		return fmt.Errorf("Cannot bind Vertex Array id 0")
	}
	gl.BindVertexArray(vao.Id)
	return nil
}

/* Draws every vertex in the array */
func (vao VertexArray) Draw() {
	vao.DrawElements(0, vao.Length)
}

/* Draws a range of verticies in the vertex array */
func (vao VertexArray) DrawElements(start, count int32) error {
	if start < 0 || start+count > vao.Length {
		return fmt.Errorf("VAO Draw index out of range")
	}

	// bind vertex array
	err := vao.Bind()
	if err != nil {
		return err
	}

	// draw call
	gl.DrawArrays(vao.Type, 0, vao.Length)
	return nil
}
