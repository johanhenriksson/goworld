package game

import (
	"encoding/gob"
	"fmt"
	"os"
)

/* Chunks are smallest renderable units of voxel geometry */
type Chunk struct {
	Seed       int
	Cx, Cz     int
	Ox, Oy, Oz int
	Sx, Sy, Sz int
	Data       Voxels
	Light      *LightVolume
}

func NewChunk(size, seed, cx, cz int) *Chunk {
	return &Chunk{
		Data:  make(Voxels, size*size*size),
		Light: NewLightVolume(size, size, size),
		Seed:  seed,
		Cx:    cx,
		Cz:    cz,
		Ox:    cx * size,
		Oz:    cz * size,
		Sx:    size,
		Sy:    size,
		Sz:    size,
	}
}

/* Clears all voxel data in this chunk */
func (c *Chunk) Clear() {
	for i := 0; i < len(c.Data); i++ {
		c.Data[i] = EmptyVoxel
	}
}

/* Returns the slice offset for a given set of coordinates, as
   well as a bool indicating whether the position is within bounds.
   If the point is out of bounds, zero is returned */
func (c *Chunk) offset(x, y, z int) (int, bool) {
	if x < 0 || x >= c.Sx || y < 0 || y >= c.Sy || z < 0 || z >= c.Sz {
		return 0, false
	}
	pos := z*c.Sx*c.Sy + x*c.Sy + y
	return pos, true
}

/* Returns a pointer to the voxel defintion at the given position.
   If the space is empty, nil is returned */
func (c *Chunk) At(x, y, z int) Voxel {
	pos, ok := c.offset(x, y, z)
	if !ok {
		return EmptyVoxel
	}
	return c.Data[pos]
}

/* Sets a voxel. If it is outside bounds, nothing happens */
func (c *Chunk) Set(x, y, z int, voxel Voxel) {
	pos, ok := c.offset(x, y, z)
	if !ok {
		return
	}
	c.Data[pos] = voxel
}

func (c *Chunk) Free(x, y, z int) bool {
	v, ok := c.offset(x, y, z)
	if !ok {
		return true
	}
	return c.Data[v] == EmptyVoxel
}

func (c *Chunk) Write(path string) error {
	filepath := fmt.Sprintf("%s/c_%d_%d.bin", path, c.Cx, c.Cz)
	file, err := os.Create(filepath)
	defer file.Close()
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(c)
	if err == nil {
		fmt.Printf("Wrote chunk %d,%d to disk\n", c.Cx, c.Cz)
	} else {
		fmt.Printf("Error writing chunk %d,%d: %s\n", c.Cx, c.Cz, err)
	}
	return err
}

func LoadChunk(path string, cx, cz int) (*Chunk, error) {
	filepath := fmt.Sprintf("%s/c_%d_%d.bin", path, cx, cz)
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	decoder := gob.NewDecoder(file)
	chunk := &Chunk{}
	err = decoder.Decode(chunk)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Read chunk %d,%d from disk\n", chunk.Cx, chunk.Cz)
	return chunk, nil
}
