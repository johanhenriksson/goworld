package terrain

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/ivec2"
)

type Tile struct {
	Map      *Map
	Position ivec2.T
	Size     int
	points   [][]Point
}

func NewTile(m *Map, position ivec2.T, size int) *Tile {
	if size < 1 {
		panic("size must be at least 1")
	}

	noiseScale := float32(1)
	heightScale := float32(20)
	height := math.NewNoise(11000, 1.0/40.0)
	noise := math.NewNoise(418941, 0.2)

	vertices := size + 1
	points := make([][]Point, vertices)
	for z := 0; z < vertices; z++ {
		points[z] = make([]Point, vertices)
		for x := 0; x < vertices; x++ {
			wx, wz := x+size*position.X, z+size*position.Y
			points[z][x].Height = heightScale*height.Sample(wx, 0, wz) + noiseScale*noise.Sample(wx, 0, wz)
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

// Patch returns a patch of the tile
func (t *Tile) Patch(offset, size ivec2.T) *Patch {
	if offset.X < 0 || offset.Y < 0 || offset.X+size.X > t.Size || offset.Y+size.Y > t.Size {
		panic("patch out of bounds")
	}
	points := make([][]Point, size.Y)
	for z := 0; z < size.Y; z++ {
		points[z] = make([]Point, size.X)
		for x := 0; x < size.X; x++ {
			points[z][x] = t.Point(offset.X+x, offset.Y+z)
		}
	}
	return &Patch{
		Size:   size,
		Offset: offset,
		Points: points,
		Source: t,
	}
}

// ApplyPatch applies a patch to the tile
func (t *Tile) ApplyPatch(p *Patch) {
	if p.Source != t {
		panic("patch source must be tile")
	}
	for z := 0; z < p.Size.Y; z++ {
		for x := 0; x < p.Size.X; x++ {
			t.SetPoint(p.Offset.X+x, p.Offset.Y+z, p.Points[z][x])
		}
	}
}
