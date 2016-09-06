package render

import (
    "fmt"
    "github.com/go-gl/gl/v4.1-core/gl"
)

type VertexArray struct {
    Id      uint32  /* OpenGL Vertex Array identifier */
    Type    uint32  /* Primitive type */
    Length  int32   /* Number of primitives */
}

func CreateVertexArray() *VertexArray {
    vao := &VertexArray {
        Type: gl.TRIANGLES,
    }
    gl.GenVertexArrays(1, &vao.Id)
    return vao
}

func (vao *VertexArray) Bind() {
    gl.BindVertexArray(vao.Id)
}

func (vao *VertexArray) Draw() {
    vao.DrawElements(0, vao.Length)
}

func (vao VertexArray) DrawElements(start, count int32) error {
    if start < 0 || start + count > vao.Length {
        // todo: error
        return fmt.Errorf("Draw index out of range")
    }

    err := vao.Bind()
    if err != nil { return err }

    gl.DrawArrays(vao.Type, 0, vao.Length)
}
