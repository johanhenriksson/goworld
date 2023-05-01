package terrain

import (
	"github.com/johanhenriksson/goworld/math/ivec2"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/render/color"
)

type Map struct {
	TileSize int
	Color    color.T

	tiles map[ivec2.T]*Tile
}

func (m *Map) getTile(tx, ty int, create bool) *Tile {
	tp := ivec2.New(tx, ty)
	tile, exists := m.tiles[tp]
	if !exists {
		if create {
			t := NewTile(tp, m.TileSize, m.Color)
			m.tiles[tp] = t
			return t
		}
	}
	return tile
}

func (m *Map) Get(point vec2.T) (Point, bool) {
	p := point.Floor()
	x, y := int(p.X), int(p.Y)

	tx, ty := x/m.TileSize, y/m.TileSize
	ox, oy := x%m.TileSize, y%m.TileSize

	tile := m.getTile(tx, ty, false)
	if tile == nil {
		return Point{}, false
	}

	return tile.points[oy][ox], true
}

func (m *Map) Set(point vec2.T, data Point) {
	p := point.Floor()
	x, y := int(p.X), int(p.Y)

	tx, ty := x/m.TileSize, y/m.TileSize
	ox, oy := x%m.TileSize, y%m.TileSize

	create := data.Height > 0 || data.R > 0 || data.G > 0 || data.B > 0
	tile := m.getTile(tx, ty, create)
	if tile == nil {
		return
	}

	tile.points[oy][ox] = data

	// if its an edge point, update neighbors accordingly
	mt := m.TileSize - 1
	if ox == 0 {
		nb := m.getTile(tx-1, ty, false)
		if nb != nil {
			nb.points[oy][mt] = data
		}
	}
	if oy == 0 {
		nb := m.getTile(tx, ty-1, false)
		if nb != nil {
			nb.points[mt][ox] = data
		}
	}
	if ox == mt {
		nb := m.getTile(tx+1, ty, false)
		if nb != nil {
			nb.points[oy][0] = data
		}
	}
	if oy == mt {
		nb := m.getTile(tx, ty+1, false)
		if nb != nil {
			nb.points[0][ox] = data
		}
	}
}
