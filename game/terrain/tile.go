package terrain

import (
	"github.com/johanhenriksson/goworld/core/events"
	"github.com/johanhenriksson/goworld/math/ivec2"
	"github.com/johanhenriksson/goworld/render/texture"
)

type Tile struct {
	Position ivec2.T
	Size     int
	UVScale  float32
	Textures []texture.Ref
	Changed  events.Event[*Tile]

	points [][]Point
}

func NewTile(position ivec2.T, size int) *Tile {
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

	// scale textures so that 1 unit = 1 pixel
	// todo: aquire texture size from texture
	textureSize := 16
	uvscale := 2 * float32(size) / float32(textureSize)

	return &Tile{
		Position: position,
		Size:     size,
		UVScale:  uvscale,
		Textures: []texture.Ref{
			texture.PathArgsRef("textures/terrain/grass1.png", texture.Args{
				Filter: texture.FilterNearest,
			}),
			texture.PathArgsRef("textures/terrain/grass2.png", texture.Args{
				Filter: texture.FilterNearest,
			}),
			texture.PathArgsRef("textures/terrain/path1.png", texture.Args{
				Filter: texture.FilterNearest,
			}),
			texture.PathArgsRef("textures/terrain/sand1.png", texture.Args{
				Filter: texture.FilterNearest,
			}),
		},

		points: points,
	}
}

func (t *Tile) Point(x, z int) Point {
	return t.points[z][x]
}

func (t *Tile) SetPoint(x, z int, p Point) {
	t.points[z][x] = p
}
