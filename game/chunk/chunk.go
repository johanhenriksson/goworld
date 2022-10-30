package chunk

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/johanhenriksson/goworld/game/voxel"
)

// Chunk is the smallest individually renderable unit of voxel geometry
type T struct {
	Seed       int
	Cx, Cz     int
	Ox, Oy, Oz int
	Sx, Sy, Sz int
	Data       voxel.Array
	Light      *LightVolume
}

func New(size, seed, cx, cz int) *T {
	return &T{
		Data:  make(voxel.Array, size*size*size),
		Light: NewLightVolume(size, size+1, size),
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

// Clear all voxel data in this chunk
func (c *T) Clear() {
	for i := 0; i < len(c.Data); i++ {
		c.Data[i] = voxel.Empty
	}
	c.Light.Clear()
}

/*
Returns the slice offset for a given set of coordinates, as

	well as a bool indicating whether the position is within bounds.
	If the point is out of bounds, zero is returned
*/
func (c *T) offset(x, y, z int) (int, bool) {
	if x < 0 || x >= c.Sx || y < 0 || y >= c.Sy || z < 0 || z >= c.Sz {
		return 0, false
	}
	pos := z*c.Sx*c.Sy + x*c.Sy + y
	return pos, true
}

/*
Returns a pointer to the voxel defintion at the given position.

	If the space is empty, nil is returned
*/
func (c *T) At(x, y, z int) voxel.T {
	pos, ok := c.offset(x, y, z)
	if !ok {
		return voxel.Empty
	}
	return c.Data[pos]
}

// Set a voxel. If it's out of bounds, nothing happens
func (c *T) Set(x, y, z int, vox voxel.T) {
	pos, ok := c.offset(x, y, z)
	if !ok {
		return
	}
	c.Data[pos] = vox
	c.Light.Block(x, y, z, vox != voxel.Empty)
}

// Free returns true if the given position is open
func (c *T) Free(x, y, z int) bool {
	v, ok := c.offset(x, y, z)
	if !ok {
		return true
	}
	c.Light.Block(x, y, z, false)
	return c.Data[v] == voxel.Empty
}

// Write the chunk to disk
func (c *T) Write(path string) error {
	filepath := fmt.Sprintf("%s/c_%d_%d.bin", path, c.Cx, c.Cz)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(c)
	if err == nil {
		fmt.Printf("Wrote chunk %d,%d to disk\n", c.Cx, c.Cz)
	} else {
		fmt.Printf("Error writing chunk %d,%d: %s\n", c.Cx, c.Cz, err)
	}
	return err
}

func Load(path string, cx, cz int) (*T, error) {
	filepath := fmt.Sprintf("%s/c_%d_%d.bin", path, cx, cz)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	chunk := &T{}
	err = decoder.Decode(chunk)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Read chunk %d,%d from disk\n", chunk.Cx, chunk.Cz)
	return chunk, nil
}
