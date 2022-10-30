package voxel

import (
	"github.com/johanhenriksson/goworld/render/color"
)

// voxel.Empty is an empty color voxel
var Empty = T{}

// Array is a collection of voxels
type Array []T

// T holds color information for a single colored voxel
type T struct {
	R, G, B byte
}

// New creates a new Color Voxel from a given color
func New(color color.T) T {
	return T{
		R: byte(255 * color.R),
		G: byte(255 * color.G),
		B: byte(255 * color.B),
	}
}
