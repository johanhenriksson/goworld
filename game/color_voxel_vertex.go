package game

// ColorVoxelVertex represents a single RGB-colored voxel
type ColorVoxelVertex struct {
	X, Y, Z byte // position
	N       byte // normal index
	R, G, B byte // color
	O       byte // occlusion
}

// ColorVoxelVertices holds a collection of vertices. Satisfies the Vertex Data interface
type ColorVoxelVertices []ColorVoxelVertex

func (buffer ColorVoxelVertices) Elements() int {
	return len(buffer)
}

func (buffer ColorVoxelVertices) Size() int {
	return 8
}
