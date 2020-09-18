package game

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
)

/* Chunks are smallest renderable units of voxel geometry */
type ColorChunk struct {
	*engine.Mesh

	Size int
	Seed int
	Ox   int
	Oy   int
	Oz   int
	Data ColorVoxels
}

func NewColorChunk(parentObject *engine.Object, size int) *ColorChunk {
	mesh := engine.NewMesh(assets.GetMaterialCached("color_voxels"))

	chk := &ColorChunk{
		Mesh: mesh,
		Size: size,
		Data: make(ColorVoxels, size*size*size),
	}

	chk.ComponentBase = engine.NewComponent(parentObject, chk)
	return chk
}

/* Clears all voxel data in this chunk */
func (chk *ColorChunk) Clear() {
	for i := 0; i < len(chk.Data); i++ {
		chk.Data[i] = EmptyColorVoxel
	}
}

/* Returns the slice offset for a given set of coordinates, as
   well as a bool indicating whether the position is within bounds.
   If the point is out of bounds, zero is returned */
func (chk *ColorChunk) offset(x, y, z int) (int, bool) {
	s := chk.Size
	if x < 0 || x >= s || y < 0 || y >= s || z < 0 || z >= s {
		return 0, false
	}
	s2 := s * s
	pos := z*s2 + y*s + x
	return pos, true
}

/* Returns a pointer to the voxel defintion at the given position.
   If the space is empty, nil is returned */
func (chk *ColorChunk) At(x, y, z int) ColorVoxel {
	pos, ok := chk.offset(x, y, z)
	if !ok {
		return EmptyColorVoxel
	}
	return chk.Data[pos]
}

/* Sets a voxel. If it is outside bounds, nothing happens */
func (chk *ColorChunk) Set(x, y, z int, voxel ColorVoxel) {
	pos, ok := chk.offset(x, y, z)
	if !ok {
		return
	}
	chk.Data[pos] = voxel
}

func (chk *ColorChunk) Free(x, y, z int) bool {
	v, ok := chk.offset(x, y, z)
	if !ok {
		return true
	}
	return chk.Data[v] == EmptyColorVoxel
}

// Computes the chunk mesh and returns a pointer to it.
func (chk *ColorChunk) Compute() {
	s := chk.Size
	data := make(ColorVoxelVertices, 0, 64)

	light := NewLightVolume(chk.Size)
	for z := 0; z < light.Size; z++ {
		for y := 0; y < light.Size; y++ {
			for x := 0; x < light.Size; x++ {
				if !chk.Free(x, y, z) {
					light.Block(x, y, z)
				}
			}
		}
	}
	light.Calculate()

	/* geometry pass */
	for z := 0; z < s; z++ {
		for y := 0; y < s; y++ {
			for x := 0; x < s; x++ {
				v := chk.At(x, y, z)
				if v == EmptyColorVoxel {
					continue
				}

				/* Simple optimization - dont draw hidden faces */
				xp := chk.At(x+1, y, z) == EmptyColorVoxel
				xn := chk.At(x-1, y, z) == EmptyColorVoxel
				yp := chk.At(x, y+1, z) == EmptyColorVoxel
				yn := chk.At(x, y-1, z) == EmptyColorVoxel
				zp := chk.At(x, y, z+1) == EmptyColorVoxel
				zn := chk.At(x, y, z-1) == EmptyColorVoxel

				/* Compute & append vertex data */
				vertices := v.Compute(light, byte(x), byte(y), byte(z), xp, xn, yp, yn, zp, zn)
				data = append(data, vertices...)
			}
		}
	}

	// buffer vertex data to GPU memory
	chk.Mesh.Buffer("geometry", data)
}
