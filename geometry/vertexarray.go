package geometry

import (
    "github.com/go-gl/gl/v4.1-core/gl"
)

type VertexArray struct {
    Id      uint32
    Length  int32
    Type    uint32
}

func CreateVertexArray() *VertexArray {
    var vao uint32
    gl.GenVertexArrays(1, &vao)
    return &VertexArray {
        Id: vao,
        Type: gl.TRIANGLES,
    }
}

func (vao *VertexArray) Bind() {
    gl.BindVertexArray(vao.Id)
}

func (vao *VertexArray) Draw() {
    vao.Bind()
    gl.DrawArrays(vao.Type, 0, vao.Length)
}
