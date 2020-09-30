package game

import (
	"github.com/johanhenriksson/goworld/math"
)

type WorldGenerator struct {
	Seed     int
	Size     int
	Rock     *math.Noise
	Grass    *math.Noise
	Cave     *math.Noise
	Variance *math.Noise
}

func ExampleWorldgen(seed, size int) *WorldGenerator {
	return &WorldGenerator{
		Seed:     seed,
		Size:     size,
		Rock:     math.NewNoise(seed+10000, 1.0/40.0),
		Grass:    math.NewNoise(seed+10002, 1.0/28.0),
		Cave:     math.NewNoise(seed+18002, 1.0/14.0),
		Variance: math.NewNoise(seed+12004, 1.0/0.5),
	}
}

func (wg *WorldGenerator) Chunk(cx, cz int) *Chunk {
	chunk := NewChunk(wg.Size, wg.Seed, cx, cz)
	for z := 0; z < chunk.Sz; z++ {
		for y := 0; y < chunk.Sy; y++ {
			for x := 0; x < chunk.Sx; x++ {
				voxel := wg.Voxel(chunk.Ox+x, chunk.Oy+y, chunk.Oz+z)
				chunk.Set(x, y, z, voxel)
				if voxel != EmptyVoxel {
					chunk.Light.Block(x, y, z, true)
				}
			}
		}
	}
	chunk.Light.Calculate()
	go chunk.Write("chunks")
	return chunk
}

func (wg *WorldGenerator) Voxel(x, y, z int) Voxel {
	rock2 := Voxel{R: 137, G: 131, B: 119}
	rock := Voxel{R: 173, G: 169, B: 158}
	grass := Voxel{R: 72, G: 140, B: 54}

	gh := int(9 * wg.Grass.Sample(x, y, z))
	rh := int(44 * wg.Rock.Sample(x, y, z))
	grassHeight := 8

	var vtype Voxel
	if y < grassHeight {
		vtype = rock2
	}

	if y == grassHeight {
		vtype = grass
	}
	if y <= grassHeight+gh && y > grassHeight {
		vtype = grass
	}
	if y < rh {
		vtype = rock
	}

	if wg.Cave.Sample(x, y, z) > 0.5 {
		vtype = EmptyVoxel
	}

	if vtype != EmptyVoxel {
		l := 1 - 0.3 * (wg.Variance.Sample(x, y, z) + 1)/2
		vtype.R = byte(l * float32(vtype.R))
		vtype.G = byte(l * float32(vtype.G))
		vtype.B = byte(l * float32(vtype.B))
	}

	return vtype
}
