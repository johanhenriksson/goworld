package terrain

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/ivec2"
	"github.com/johanhenriksson/goworld/render/color"
)

type Tile struct {
	Map      *Map
	Position ivec2.T
	Size     int
	points   [][]Point
}

func NewTile(m *Map, position ivec2.T, size int, color color.T) *Tile {
	if size < 1 {
		panic("size must be at least 1")
	}
	noise := math.NewNoise(10000, 1.0/40.0)

	vertices := size + 1
	points := make([][]Point, vertices)
	for z := 0; z < vertices; z++ {
		points[z] = make([]Point, vertices)
		for x := 0; x < vertices; x++ {
			b := color.Byte4()
			points[z][x].R = b.X
			points[z][x].G = b.Y
			points[z][x].B = b.Z

			wx, wz := x+size*position.X, z+size*position.Y
			points[z][x].Height = 40 * noise.Sample(wx, 0, wz)
		}
	}

	return &Tile{
		Map:      m,
		Position: position,
		Size:     size,
		points:   points,
	}
}

func (t *Tile) Point(x, z int) Point {
	return t.points[z][x]
}
