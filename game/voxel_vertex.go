package game

import "unsafe"

// VoxelVertex represents a single RGB-colored voxel
type VoxelVertex struct {
	X, Y, Z byte // position
	N       byte // normal index
	R, G, B byte // color
	O       byte // occlusion
}

// VoxelVertices holds a collection of vertices. Satisfies the Vertex Data interface
type VoxelVertices []VoxelVertex

func (buffer VoxelVertices) Elements() int {
	return len(buffer)
}

func (buffer VoxelVertices) Size() int {
	return 8
}

func (buffer VoxelVertices) Pointer() unsafe.Pointer {
	return unsafe.Pointer(&buffer[0])
}
