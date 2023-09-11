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

	vertices := size + 1
	points := make([][]Point, vertices)
	for z := 0; z < vertices; z++ {
		points[z] = make([]Point, vertices)
		for x := 0; x < vertices; x++ {
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
