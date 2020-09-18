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

func NewLightVolume(size int) *LightVolume {
	lvox := make([][][]LightVoxel, size)
	for z := 0; z < size; z++ {
		lvox[z] = make([][]LightVoxel, size)
		for y := 0; y < size; y++ {
			lvox[z][y] = make([]LightVoxel, size)
		}
	}
	return &LightVolume{
		Falloff: 0.7,
		Size:    size,
		data:    lvox,
	}
}

func (lv *LightVolume) reset() {
	for z := 0; z < lv.Size; z++ {
		for y := 0; y < lv.Size; y++ {
			for x := 0; x < lv.Size; x++ {
				lv.data[z][y][x].Calculated = false
			}
		}
	}
}

func (lv *LightVolume) Brightness(x, y, z int) float32 {
	if x < 0 || y < 0 || z < 0 || x >= lv.Size || y >= lv.Size || z >= lv.Size {
		return 1
	}
	return lv.data[z][y][x].V
}

func (lv *LightVolume) Block(x, y, z int) {
	lv.data[z][y][x].Blocked = true
}

func (lv *LightVolume) Calculate() {
	lv.reset()
	i := 0
	for {
		i++
		changed := false
		for y := 0; y < lv.Size; y++ {
			for z := 0; z < lv.Size; z++ {
				for x := 0; x < lv.Size; x++ {
					lp := &lv.data[z][y][x]
					if lp.Calculated {
						continue
					}
					if lp.Blocked {
						lp.V = 0
						lp.Calculated = true
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
						lp.Calculated = true
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

type LightVoxel struct {
	Blocked    bool
	Calculated bool
	V          float32
}

func (lp LightVoxel) String() string {
	if lp.Blocked {
		return "    "
	} else {
		return fmt.Sprintf("%.2f", lp.V)
	}
}
