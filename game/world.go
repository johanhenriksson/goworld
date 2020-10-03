package game

import (
	"fmt"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type ChunkProvider interface {
	Chunk(x, z int) *Chunk
	Voxel(x, y, z int) Voxel
}

type ChunkPos struct {
	X int
	Z int
}

type World struct {
	*engine.Object

	Seed         int
	ChunkSize    int
	KeepDistance int
	DrawDistance int
	Cache        map[ChunkPos]*Chunk
	Provider     ChunkProvider
}

func NewWorld(seed, size int) *World {
	return &World{
		Object:       engine.NewObject(vec3.Zero, vec3.Zero),
		Seed:         seed,
		KeepDistance: 5,
		DrawDistance: 3,
		ChunkSize:    size,
		Cache:        make(map[ChunkPos]*Chunk),
		Provider:     ExampleWorldgen(seed, size),
	}
}

func (w *World) LoadAround(pos mgl.Vec3) {

}

func (w *World) AddChunk(cx, cz int) *Chunk {
	chunk, err := LoadChunk("chunks", cx, cz)
	if err != nil {
		chunk = w.Provider.Chunk(cx, cz)
		fmt.Printf("Generated chunk %d,%d\n", cx, cz)
	}

	w.Cache[ChunkPos{cx, cz}] = chunk
	return chunk
}

func (w *World) Voxel(x, y, z int) Voxel {
	cx, cz := x/w.ChunkSize, z/w.ChunkSize
	lx, ly, lz := x%w.ChunkSize, y, z%w.ChunkSize
	if chunk, exists := w.Cache[ChunkPos{cx, cz}]; exists {
		return chunk.At(lx, ly, lz)
	}
	return w.Provider.Voxel(x, y, z)
}

func (w *World) Set(x, y, z int, voxel Voxel) {
	cx, cz := x/w.ChunkSize, z/w.ChunkSize
	lx, ly, lz := x%w.ChunkSize, y, z%w.ChunkSize
	if chunk, exists := w.Cache[ChunkPos{cx, cz}]; exists {
		chunk.Set(lx, ly, lz, voxel)
		chunk.Light.Block(lx, ly, lz, voxel != EmptyVoxel)
	}
}

func (w *World) HeightAt(p vec3.T) float32 {
	x, y, z := int(p.X), int(p.Y), int(p.Z)
	for w.Voxel(x, y, z) == EmptyVoxel && y >= 0 {
		y--
	}
	y++
	return float32(y)
}
