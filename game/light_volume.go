package game

import (
	"fmt"
	"github.com/johanhenriksson/goworld/math"
)

type LightVolume struct {
	Falloff    float32
	Sx, Sy, Sz int
	data       [][][]LightVoxel
}

type LightVoxel struct {
	Blocked bool
	V       float32
}

func NewLightVolume(sx, sy, sz int) *LightVolume {
	lvox := make([][][]LightVoxel, sz)
	for z := 0; z < sz; z++ {
		lvox[z] = make([][]LightVoxel, sx)
		for x := 0; x < sx; x++ {
			lvox[z][x] = make([]LightVoxel, sy)
		}
	}
	return &LightVolume{
		Falloff: 0.6,
		Sx:      sx,
		Sy:      sy,
		Sz:      sz,
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
	if x < 0 || y < 0 || z < 0 || x >= lv.Sx || y >= lv.Sy || z >= lv.Sz {
		return nil
	}
	return &lv.data[z][x][y]
}

func (lv *LightVolume) Calculate() {
	i := 0
	for {
		i++
		changed := false
		for z := 0; z < lv.Sz; z++ {
			for x := 0; x < lv.Sx; x++ {
				for y := 0; y < lv.Sy; y++ {
					lp := lv.Get(x, y, z)
					if lp.Blocked {
						lp.V = 0
						continue
					}

					nmax := float32(0)
					if x > 0 {
						nmax = math.Max(nmax, lv.Get(x-1, y, z).V*lv.Falloff)
					}
					if y > 0 {
						nmax = math.Max(nmax, lv.Get(x, y-1, z).V*lv.Falloff)
					}
					if z > 0 {
						nmax = math.Max(nmax, lv.Get(x, y, z-1).V*lv.Falloff)
					}
					if x < lv.Sx-1 {
						nmax = math.Max(nmax, lv.Get(x+1, y, z).V*lv.Falloff)
					}
					if y < lv.Sy-1 {
						nmax = math.Max(nmax, lv.Get(x, y+1, z).V)
					}
					if z < lv.Sz-1 {
						nmax = math.Max(nmax, lv.Get(x, y, z+1).V*lv.Falloff)
					}

					if y == lv.Sy-1 {
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
