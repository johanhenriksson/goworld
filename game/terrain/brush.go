package terrain

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/ivec2"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// A brush is a tool for modifying terrain
type Brush interface {
	// Paint applies the brush to a patch of terrain
	Paint(patch *Patch, center vec3.T, radius, dt float32) error
}

// A patch is a rectangular area of a tile
type Patch struct {
	// Size is the size of the patch
	Size ivec2.T

	// Offset is the offset of the patch in the tile
	Offset ivec2.T

	// Points holds the terrain data for the patch. Ordered as [z][x].
	// The points can be modified without affecting the original tile
	Points [][]Point

	// Source is the tile that the patch was copied from
	Source *Tile
}

// A raise brush raises or lowers the terrain
type RaiseBrush struct{}

func (b *RaiseBrush) Paint(patch *Patch, center vec3.T, radius, strength float32) error {
	// implement brush operation
	// apply operations on copied points
	for z := 0; z < patch.Size.Y; z++ {
		for x := 0; x < patch.Size.X; x++ {
			p := vec2.NewI(patch.Offset.X+x, patch.Offset.Y+z)

			// calculate brush weight as the distance from center of brush
			weight := math.Min(p.Sub(center.XZ()).Length()/radius, 1)

			// invert
			weight = 1 - weight

			// quadratic falloff
			weight = weight * weight

			patch.Points[z][x].Height += strength * weight
		}
	}

	return nil
}
