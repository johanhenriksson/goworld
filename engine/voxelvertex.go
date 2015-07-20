package engine

import (
    "unsafe"
    "github.com/go-gl/gl/v4.1-core/gl"
)

type VoxelVertex struct {
    X, Y, Z     uint8 // Vertex position relative to chunk
    Nx, Ny, Nz  int8  // Normal vector
    Tx, Ty      uint8 // Tile tex coords
}

type VoxelVertices []VoxelVertex

func (buffer VoxelVertices) Elements() int {
    return len(buffer)
}

func (buffer VoxelVertices) Size() int {
    return 8
}

func (buffer VoxelVertices) GLPtr() unsafe.Pointer {
    return gl.Ptr(buffer)
}

