package game

// VoxelVertex represents a single RGB-colored voxel
type VoxelVertex struct {
	X byte `vtx:"position,uint8,3"`
	Y byte `vtx:"skip"`
	Z byte `vtx:"skip"`
	N byte `vtx:"normal_id,uint8,1"`
	R byte `vtx:"color,uint8,3,normalize"`
	G byte `vtx:"skip"`
	B byte `vtx:"skip"`
	O byte `vtx:"occlusion,uint8,1,normalize"`
}
