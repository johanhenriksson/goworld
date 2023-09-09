package terrain

import (
	"sync"

	"github.com/johanhenriksson/goworld/math/ivec2"
	"github.com/johanhenriksson/goworld/math/vec2"
)

type Map struct {
	TileSize int

	tiles map[ivec2.T]*Tile
	mutex sync.Mutex
}

func NewMap(tileSize int, tiles int) *Map {
	m := &Map{
		TileSize: tileSize,
		tiles:    make(map[ivec2.T]*Tile),
	}
	for z := 0; z < tiles; z++ {
		for x := 0; x < tiles; x++ {
			tp := ivec2.New(x, z)
			m.tiles[tp] = NewTile(m, tp, m.TileSize)
		}
	}
	return m
}

func (m *Map) GetTile(tx, ty int, create bool) *Tile {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tp := ivec2.New(tx, ty)
	tile, exists := m.tiles[tp]
	if !exists {
		if create {
			t := NewTile(m, tp, m.TileSize)
			m.tiles[tp] = t
			return t
		} else {
			return nil
		}
	}
	return tile
}

func (m *Map) Get(point vec2.T) (Point, bool) {
	p := point.Floor()
	x, y := int(p.X), int(p.Y)

	tx, ty := x/m.TileSize, y/m.TileSize
	ox, oy := x%m.TileSize, y%m.TileSize

	tile := m.GetTile(tx, ty, false)
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

	create := data.Height != 0 || data.Weights[0] > 0 || data.Weights[1] > 0 || data.Weights[2] > 0 || data.Weights[3] > 0
	tile := m.GetTile(tx, ty, create)
	if tile == nil {
		return
	}

	tile.points[oy][ox] = data

	// if its an edge point, update neighbors accordingly
	mt := m.TileSize - 1
	if ox == 0 {
		nb := m.GetTile(tx-1, ty, false)
		if nb != nil {
			nb.points[oy][mt] = data
		}
	}
	if oy == 0 {
		nb := m.GetTile(tx, ty-1, false)
		if nb != nil {
			nb.points[mt][ox] = data
		}
	}
	if ox == mt {
		nb := m.GetTile(tx+1, ty, false)
		if nb != nil {
			nb.points[oy][0] = data
		}
	}
	if oy == mt {
		nb := m.GetTile(tx, ty+1, false)
		if nb != nil {
			nb.points[0][ox] = data
		}
	}
}
