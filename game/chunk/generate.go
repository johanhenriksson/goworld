package chunk

import "github.com/johanhenriksson/goworld/game/voxel"

type Generator interface {
	Voxel(x, y, z int) voxel.T
}

func Generate(wg Generator, size, cx, cz int) *T {
	chonk := New(size, cx, cz)
	for z := 0; z < chonk.Sz; z++ {
		for y := 0; y < chonk.Sy; y++ {
			for x := 0; x < chonk.Sx; x++ {
				vox := wg.Voxel(chonk.Ox+x, chonk.Oy+y, chonk.Oz+z)
				chonk.Set(x, y, z, vox)
				if vox != voxel.Empty {
					chonk.Light.Block(x, y, z, true)
				}
			}
		}
	}
	chonk.Light.Calculate()
	go chonk.Write("chunks")
	return chonk
}
