package chunk

import (
	"math/rand"

	"github.com/johanhenriksson/goworld/game/voxel"
	"github.com/johanhenriksson/goworld/render/color"
)

func NewRandomGen() Generator {
	return &randomGen{}
}

type randomGen struct{}

func (rg *randomGen) Voxel(x, y, z int) voxel.T {
	i := rand.Intn(len(color.DefaultPalette))
	c := color.DefaultPalette[i].Byte4()
	return voxel.T{
		R: c.X,
		G: c.Y,
		B: c.Z,
	}
}
