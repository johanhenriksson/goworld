package game

import (
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/engine"
)

type ChunkProvider interface {
	Get(x, y, z int) Voxels
	Update(x, y, z int, v Voxel)
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
		Object:       engine.NewObject(0, 0, 0),
		Seed:         seed,
		KeepDistance: 5,
		DrawDistance: 3,
		ChunkSize:    size,
		Cache:        make(map[ChunkPos]*Chunk),
	}
}

func (w *World) LoadAround(pos mgl.Vec3) {

}

func (w *World) AddChunk(cx, cz int) *Chunk {
	chunk, err := LoadChunk("chunks", cx, cz)
	if err != nil {
		chunk = NewChunk(w.ChunkSize, w.Seed, cx, cz)
		// generate voxel data
		GenerateChunk(chunk)
	}
	//obj := engine.NewObject(float32(cx*w.ChunkSize), 0, float32(cz*w.ChunkSize))

	w.Cache[ChunkPos{X: cx, Z: cz}] = chunk
	return chunk
}
