package chunk

import (
	"fmt"

	"github.com/johanhenriksson/goworld/game/voxel"
)

type Generator interface {
	Voxel(x, y, z int) voxel.T
}

func Generate(wg Generator, size, cx, cz int) *T {
	offsetX, offsetZ := cx*size, cz*size
	chonk := New(fmt.Sprintf("%d_%d", cx, cz), size, size, size)
	for z := 0; z < chonk.Sz; z++ {
		for y := 0; y < chonk.Sy; y++ {
			for x := 0; x < chonk.Sx; x++ {
				vox := wg.Voxel(offsetX+x, y, offsetZ+z)
				chonk.Set(x, y, z, vox)
				if vox != voxel.Empty {
					chonk.Light.Block(x, y, z, true)
				}
			}
		}
	}
	chonk.Light.Calculate()
	return chonk
}
