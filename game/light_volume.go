package game

import (
	"fmt"
	"github.com/johanhenriksson/goworld/math"
)

type LightVolume struct {
	Falloff float32
	Size    int
	data    [][][]LightVoxel
}

type LightVoxel struct {
	Blocked bool
	V       float32
}

func NewLightVolume(size int) *LightVolume {
	lvox := make([][][]LightVoxel, size)
	for z := 0; z < size; z++ {
		lvox[z] = make([][]LightVoxel, size)
		for y := 0; y < size; y++ {
			lvox[z][y] = make([]LightVoxel, size)
		}
	}
	return &LightVolume{
		Falloff: 0.6,
		Size:    size,
		data:    lvox,
	}
}

func (lv *LightVolume) Brightness(x, y, z int) float32 {
	v := lv.Get(x, y, z)
	if v != nil {
		return v.V
	}
	return 1
}

func (lv *LightVolume) Block(x, y, z int) {
	v := lv.Get(x, y, z)
	if v != nil {
		v.Blocked = true
	}
}

func (lv *LightVolume) Get(x, y, z int) *LightVoxel {
	if x < 0 || y < 0 || z < 0 || x >= lv.Size || y >= lv.Size || z >= lv.Size {
		return nil
	}
	return &lv.data[z][y][x]
}

func (lv *LightVolume) Calculate() {
	i := 0
	for {
		i++
		changed := false
		for y := 0; y < lv.Size; y++ {
			for z := 0; z < lv.Size; z++ {
				for x := 0; x < lv.Size; x++ {
					lp := &lv.data[z][y][x]
					if lp.Blocked {
						lp.V = 0
						continue
					}

					nmax := float32(0)

					if x > 0 {
						nmax = math.Max(nmax, lv.data[z][y][x-1].V*lv.Falloff)
					}
					if y > 0 {
						nmax = math.Max(nmax, lv.data[z][y-1][x].V*lv.Falloff)
					}
					if z > 0 {
						nmax = math.Max(nmax, lv.data[z-1][y][x].V*lv.Falloff)
					}
					if x < lv.Size-1 {
						nmax = math.Max(nmax, lv.data[z][y][x+1].V*lv.Falloff)
					}
					if y < lv.Size-1 {
						nmax = math.Max(nmax, lv.data[z][y+1][x].V)
					}
					if z < lv.Size-1 {
						nmax = math.Max(nmax, lv.data[z+1][y][x].V*lv.Falloff)
					}

					if y == lv.Size-1 {
						nmax = 1
					}

					if lp.V != nmax {
						lp.V = nmax
						changed = true
					}
				}
			}
		}
		if !changed {
			fmt.Println("Light volume calculation finished in", i, "iterations")
			break
		}
	}
}
