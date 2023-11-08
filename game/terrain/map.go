package terrain

import (
	"sync"

	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/ivec2"
)

type Map struct {
	TileSize int
	Name     string

	tiles map[ivec2.T]*Tile
	mutex sync.Mutex
}

func NewMap(tileSize int) *Map {
	m := &Map{
		TileSize: tileSize,
		Name:     "default",
		tiles:    make(map[ivec2.T]*Tile),
	}
	return m
}

func (m *Map) Tile(tx, ty int, create bool) *Tile {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	tp := ivec2.New(tx, ty)
	tile, exists := m.tiles[tp]
	if !exists {
		if create {
			t := NewTile(tp, m.TileSize)
			m.tiles[tp] = t
			return t
		} else {
			return nil
		}
	}
	return tile
}

func (m *Map) AddTile(t *Tile) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.tiles[t.Position] = t
}

func (m *Map) Get(offset, size ivec2.T) *Patch {
	// allocate patch
	points := make([][]Point, size.Y+1)
	for z := 0; z <= size.Y; z++ {
		points[z] = make([]Point, size.X+1)
	}
	patch := &Patch{
		Size:   size,
		Offset: offset,
		Points: points,
		Source: m,
	}

	tmin := m.TileCoords(patch.Offset)
	tmax := m.TileCoords(patch.Offset.Add(patch.Size))

	for x := tmin.X; x <= tmax.X; x++ {
		for z := tmin.Y; z <= tmax.Y; z++ {
			t := m.Tile(x, z, true)
			Tx := x * m.TileSize
			Tz := z * m.TileSize

			// find the region of the patch that belongs to this tile
			tmx := math.Max(patch.Offset.X-Tx, 0)
			tmz := math.Max(patch.Offset.Y-Tz, 0)
			tMx := math.Min(patch.Offset.X+patch.Size.X-Tx, m.TileSize)
			tMz := math.Min(patch.Offset.Y+patch.Size.Y-Tz, m.TileSize)

			px := math.Max(Tx-patch.Offset.X, 0) - tmx
			pz := math.Max(Tz-patch.Offset.Y, 0) - tmz

			// copy the patch region to the tile
			for tz := tmz; tz <= tMz; tz++ {
				for tx := tmx; tx <= tMx; tx++ {
					patch.Points[tz+pz][tx+px] = t.points[tz][tx]
				}
			}
		}
	}

	return patch
}

func (m *Map) TileCoords(point ivec2.T) ivec2.T {
	return ivec2.New(Floor(point.X, m.TileSize), Floor(point.Y, m.TileSize))
}

func (m *Map) Set(patch *Patch) {
	tmin := m.TileCoords(patch.Offset)
	tmax := m.TileCoords(patch.Offset.Add(patch.Size))

	for x := tmin.X; x <= tmax.X; x++ {
		for z := tmin.Y; z <= tmax.Y; z++ {
			t := m.Tile(x, z, true)
			Tx := x * m.TileSize
			Tz := z * m.TileSize

			// find the region of the patch that belongs to this tile
			tmx := math.Max(patch.Offset.X-Tx, 0)
			tMx := math.Min(patch.Offset.X+patch.Size.X-Tx, m.TileSize)
			tmz := math.Max(patch.Offset.Y-Tz, 0)
			tMz := math.Min(patch.Offset.Y+patch.Size.Y-Tz, m.TileSize)

			px := math.Max(Tx-patch.Offset.X, 0) - tmx
			pz := math.Max(Tz-patch.Offset.Y, 0) - tmz

			// copy the patch region to the tile
			for tz := tmz; tz <= tMz; tz++ {
				for tx := tmx; tx <= tMx; tx++ {
					t.points[tz][tx] = patch.Points[tz+pz][tx+px]
				}
			}
		}
	}

	for x := tmin.X; x <= tmax.X; x++ {
		for z := tmin.Y; z <= tmax.Y; z++ {
			t := m.Tile(x, z, false)
			t.Changed.Emit(t)
		}
	}
}

func Floor(v int, s int) int {
	if v < 0 {
		return (v - s + 1) / s
	}
	return v / s
}
