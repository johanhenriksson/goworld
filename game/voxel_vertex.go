package game

import "unsafe"

// VoxelVertex represents a single RGB-colored voxel
type VoxelVertex struct {
	X byte `vtx:"position,uint8,3"`
	Y byte
	Z byte
	N byte `vtx:"normal_id,uint8,1"`
	R byte `vtx:"color,uint8,3,normalize"`
	G byte
	B byte
	O byte `vtx:"occlusion,uint8,1,normalize"`
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
