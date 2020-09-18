package game

import (
	"github.com/johanhenriksson/goworld/math"
)

func GenerateChunk(chk *Chunk) {
	ox, oy, oz := chk.Ox, chk.Oy, chk.Oz

	/* Define voxels */
	rock2 := Voxel{R: 137, G: 131, B: 119}
	rock := Voxel{R: 173, G: 169, B: 158}
	grass := Voxel{R: 72, G: 140, B: 54}
	cloud := Voxel{R: 255, G: 255, B: 255}

	/* Fill chunk with voxels */
	rockNoise := math.NewNoise(chk.Seed+10000, 1.0/40.0)
	grassNoise := math.NewNoise(chk.Seed+10002, 1.0/28.0)
	cloudNoise := math.NewNoise(chk.Seed+24511626, 1/40.0)
	// caveNoise := math.NewNoise(chk.Seed+123981, 1/70.0)

	grassHeight := 8

	for z := 0; z < chk.Sz; z++ {
		for y := 0; y < chk.Sy; y++ {
			for x := 0; x < chk.Sx; x++ {
				gh := int(9 * grassNoise.Sample(x+ox, oy, z+oz))
				rh := int(44 * rockNoise.Sample(x+ox, oy, z+oz))
				ch := int(8*cloudNoise.Sample(x+ox, y+oy, z+oz)) + 8

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

				if ch > 12 && y > 98-ch && y < 100+ch {
					vtype = cloud
				}

				chk.Set(x, y, z, vtype)
			}
		}
	}
}
