package game

import (
	"github.com/johanhenriksson/goworld/render"
)

// EmptyVoxel is an empty color voxel
var EmptyVoxel = Voxel{}

// Voxels is a collection of voxels
type Voxels []Voxel

// Voxel holds color information for a single colored voxel
type Voxel struct {
	R, G, B byte
}

// NewVoxel creates a new Color Voxel from a given color
func NewVoxel(color render.Color) Voxel {
	return Voxel{
		R: byte(255 * color.R),
		G: byte(255 * color.G),
		B: byte(255 * color.B),
	}
}
