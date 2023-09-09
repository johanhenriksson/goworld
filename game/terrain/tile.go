package terrain

import (
	// "github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/core/events"
	"github.com/johanhenriksson/goworld/math/ivec2"
)

type Tile struct {
	Map      *Map
	Position ivec2.T
	Size     int
	Changed  events.Event[*Tile]

	points [][]Point
}

func NewTile(m *Map, position ivec2.T, size int) *Tile {
	if size < 1 {
		panic("size must be at least 1")
	}

	// noiseScale := float32(1)
	// heightScale := float32(20)
	// height := math.NewNoise(11000, 1.0/40.0)
	// noise := math.NewNoise(418941, 0.2)

	vertices := size + 1
	points := make([][]Point, vertices)
	for z := 0; z < vertices; z++ {
		points[z] = make([]Point, vertices)
		for x := 0; x < vertices; x++ {
			// wx, wz := x+size*position.X, z+size*position.Y
			// points[z][x].Height = heightScale*height.Sample(wx, 0, wz) + noiseScale*noise.Sample(wx, 0, wz)
			points[z][x].Weights[0] = 255
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

func (t *Tile) SetPoint(x, z int, p Point) {
	t.points[z][x] = p
}
