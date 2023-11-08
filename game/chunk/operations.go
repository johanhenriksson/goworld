package chunk

import (
	"github.com/johanhenriksson/goworld/game/voxel"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Crop(chk *T) {
	xmin, xmax := chk.Sx, 0
	ymin, ymax := chk.Sy, 0
	zmin, zmax := chk.Sz, 0

	for z := 0; z < chk.Sz; z++ {
		for x := 0; x < chk.Sx; x++ {
			for y := 0; y < chk.Sy; y++ {
				// ignore empty space
				if chk.At(x, y, z) == voxel.Empty {
					continue
				}

				// keep track of the lowest corner of the bounding box
				xmin = min(xmin, x)
				ymin = min(ymin, y)
				zmin = min(zmin, z)

				// keep track of the highest corner of the bounding box
				xmax = max(xmax, x+1)
				ymax = max(ymax, y+1)
				zmax = max(zmax, z+1)
			}
		}
	}

	// copy chunk data
	cropped := New(chk.Key, max(xmax-xmin, 1), max(ymax-ymin, 1), max(zmax-zmin, 1))
	for z := 0; z < cropped.Sz; z++ {
		for x := 0; x < cropped.Sx; x++ {
			for y := 0; y < cropped.Sy; y++ {
				cropped.Set(x, y, z, chk.At(xmin+x, ymin+y, zmin+z))
			}
		}
	}

	// recalculate light
	cropped.Light.Calculate()
	*chk = *cropped
}

func Expand(chk *T) {
	expanded := New(chk.Key, chk.Sx+2, chk.Sy+2, chk.Sz+2)

	for z := 0; z < chk.Sz; z++ {
		for x := 0; x < chk.Sx; x++ {
			for y := 0; y < chk.Sy; y++ {
				expanded.Set(x+1, y+1, z+1, chk.At(x, y, z))
			}
		}
	}

	expanded.Light.Calculate()
	*chk = *expanded
}

func Subdivide(chk *T) {
	expanded := New(chk.Key, chk.Sx*2, chk.Sy*2, chk.Sz*2)

	for z := 0; z < chk.Sz; z++ {
		for x := 0; x < chk.Sx; x++ {
			for y := 0; y < chk.Sy; y++ {
				expanded.Set(2*x+0, 2*y+0, 2*z+0, chk.At(x, y, z))
				expanded.Set(2*x+0, 2*y+0, 2*z+1, chk.At(x, y, z))
				expanded.Set(2*x+0, 2*y+1, 2*z+0, chk.At(x, y, z))
				expanded.Set(2*x+0, 2*y+1, 2*z+1, chk.At(x, y, z))
				expanded.Set(2*x+1, 2*y+0, 2*z+0, chk.At(x, y, z))
				expanded.Set(2*x+1, 2*y+0, 2*z+1, chk.At(x, y, z))
				expanded.Set(2*x+1, 2*y+1, 2*z+0, chk.At(x, y, z))
				expanded.Set(2*x+1, 2*y+1, 2*z+1, chk.At(x, y, z))
			}
		}
	}

	expanded.Light.Calculate()
	*chk = *expanded
}
