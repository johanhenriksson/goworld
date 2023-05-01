package terrain

import (
	"github.com/johanhenriksson/goworld/math/ivec2"
	"github.com/johanhenriksson/goworld/render/color"
)

type Tile struct {
	Position ivec2.T
	Size     int
	points   [][]Point
}

func NewTile(position ivec2.T, size int, color color.T) *Tile {
	if size < 1 {
		panic("size must be at least 1")
	}
	points := make([][]Point, size+1)
	for y := 0; y < size; y++ {
		points[y] = make([]Point, size+1)
		for x := 0; x < size; x++ {
			b := color.Byte4()
			points[y][x].R = b.X
			points[y][x].G = b.Y
			points[y][x].B = b.Z
		}
	}

	return &Tile{
		Position: position,
		Size:     size,
		points:   points,
	}
}
