package terrain

import (
	"time"

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
	Source *Map
}

// A height brush raises or lowers the terrain
type HeightBrush struct {
	Sign float32
}

func NewRaiseBrush() *HeightBrush {
	return &HeightBrush{Sign: 1}
}

func NewLowerBrush() *HeightBrush {
	return &HeightBrush{Sign: -1}
}

func (b *HeightBrush) Paint(patch *Patch, center vec3.T, radius, strength float32) error {
	for z := 0; z < patch.Size.Y; z++ {
		for x := 0; x < patch.Size.X; x++ {
			p := vec2.NewI(patch.Offset.X+x, patch.Offset.Y+z)

			// calculate brush weight as the distance from center of brush
			weight := math.Min(p.Sub(center.XZ()).Length()/radius, 1)

			// invert
			weight = 1 - weight

			// quadratic falloff
			weight = weight * weight

			patch.Points[z][x].Height += b.Sign * strength * weight
		}
	}

	return nil
}

type SmoothBrush struct{}

func (b *SmoothBrush) Paint(patch *Patch, center vec3.T, radius, strength float32) error {
	source := patch.Source.Get(patch.Offset, patch.Size)

	for z := 0; z < patch.Size.Y; z++ {
		for x := 0; x < patch.Size.X; x++ {
			p := ivec2.New(x, z)

			// calculate new height value as the average of the surrounding points
			k := 1
			points := 0
			smoothed := float32(0)
			for i := -k; i <= k; i++ {
				for j := -k; j <= k; j++ {
					q := p.Add(ivec2.New(i, j))
					if q.X < 0 || q.Y < 0 || q.X >= patch.Size.X || q.Y >= patch.Size.Y {
						continue
					}

					// read directly from source map to avoid smoothing the smoothing
					pt := source.Points[q.Y][q.X]
					smoothed += pt.Height
					points++
				}
			}
			smoothed /= float32(points)

			// apply smoothing with strength
			patch.Points[z][x].Height = strength*smoothed + (1-strength)*patch.Points[z][x].Height
		}
	}

	return nil
}

type PaintBrush struct {
}

func (b *PaintBrush) Paint(patch *Patch, center vec3.T, radius, strength float32) error {
	for z := 0; z < patch.Size.Y; z++ {
		for x := 0; x < patch.Size.X; x++ {
			p := vec2.NewI(patch.Offset.X+x, patch.Offset.Y+z)

			// calculate brush weight as the distance from center of brush
			weight := math.Min(p.Sub(center.XZ()).Length()/radius, 1)

			// invert
			weight = 1 - weight

			patch.Points[z][x].Weights[1] += 1 // strength * weight
		}
	}
	return nil
}

type LevelBrush struct{}

func (b *LevelBrush) Paint(patch *Patch, center vec3.T, radius, strength float32) error {
	// grab center height
	cx, gz := patch.Size.X/2, patch.Size.Y/2
	desiredHeight := patch.Points[gz][cx].Height

	for z := 0; z < patch.Size.Y; z++ {
		for x := 0; x < patch.Size.X; x++ {
			p := vec2.NewI(patch.Offset.X+x, patch.Offset.Y+z)

			// calculate brush weight as the distance from center of brush
			weight := math.Min(p.Sub(center.XZ()).Length()/radius, 1)

			// invert
			weight = 1 - weight

			height := patch.Points[z][x].Height

			adjustment := (desiredHeight - height) * weight * strength

			patch.Points[z][x].Height += adjustment
		}
	}

	return nil
}

type NoiseBrush struct {
	noise *math.Noise
}

func NewNoiseBrush() *NoiseBrush {
	return &NoiseBrush{
		noise: math.NewNoise(int(time.Now().UnixNano()), 4),
	}
}

func (b *NoiseBrush) Paint(patch *Patch, center vec3.T, radius, strength float32) error {
	for z := 0; z < patch.Size.Y; z++ {
		for x := 0; x < patch.Size.X; x++ {
			p := vec2.NewI(patch.Offset.X+x, patch.Offset.Y+z)

			// calculate brush weight as the distance from center of brush
			weight := math.Min(p.Sub(center.XZ()).Length()/radius, 1)

			// invert
			weight = 1 - weight

			patch.Points[z][x].Height += b.noise.Sample(int(p.X), 0, int(p.Y)) * strength * weight
		}
	}

	return nil
}
