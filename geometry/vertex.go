package geometry

import (
    "unsafe"
    "github.com/go-gl/gl/v4.1-core/gl"
)

type Vertices []Vertex

type Vertex struct {
    X, Y, Z     float32
    R, G, B     float32
}


func (vtx Vertices) Elements() int {
    return len(vtx)
}

func (vtx Vertices) Size() int {
    return 24
}

func (vtx Vertices) GLPtr() unsafe.Pointer {
    return gl.Ptr(vtx)
}
