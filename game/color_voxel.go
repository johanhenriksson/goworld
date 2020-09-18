package game

import (
	"github.com/johanhenriksson/goworld/render"
)

// EmptyColorVoxel is an empty color voxel
var EmptyColorVoxel = ColorVoxel{}

// ColorVoxels is a collection of voxels
type ColorVoxels []ColorVoxel

// ColorVoxel holds color information for a single colored voxel
type ColorVoxel struct {
	R, G, B byte
}

// NewColorVoxel creates a new Color Voxel from a given color
func NewColorVoxel(color render.Color) ColorVoxel {
	return ColorVoxel{
		R: byte(255 * color.R),
		G: byte(255 * color.G),
		B: byte(255 * color.B),
	}
}

// Compute the verticies of this Color Voxel
func (v ColorVoxel) Compute(light *LightVolume, x, y, z byte, xp, xn, yp, yn, zp, zn bool) ColorVoxelVertices {
	data := make(ColorVoxelVertices, 0, 36)
	brightness := func(x, y, z byte) byte {
		lv := 1 - light.Brightness(int(x), int(y), int(z))
		return byte(255 * lv)
	}

	// Right (X+) N=1
	if xp {
		o := brightness(x+1, y, z)
		data = append(data,
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 1, R: v.R, G: v.G, B: v.B, O: o}, // front bottom right
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 1, R: v.R, G: v.G, B: v.B, O: o}, // back bottom right
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 1, R: v.R, G: v.G, B: v.B, O: o}, // back top right
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 1, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 1, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, N: 1, R: v.R, G: v.G, B: v.B, O: o})
	}

	// Left faces (X-) N=2
	if xn {
		o := brightness(x-1, y, z)
		data = append(data,
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 2, R: v.R, G: v.G, B: v.B, O: o}, // bottom left back
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 2, R: v.R, G: v.G, B: v.B, O: o}, // top left front
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, N: 2, R: v.R, G: v.G, B: v.B, O: o}, // bottom left front
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 2, R: v.R, G: v.G, B: v.B, O: o}, // bottom left back
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 2, R: v.R, G: v.G, B: v.B, O: o}, // top left back
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 2, R: v.R, G: v.G, B: v.B, O: o}) // top left front
	}

	// Top faces (Y+) N=3
	if yp {
		o := brightness(x, y+1, z)
		data = append(data,
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 3, R: v.R, G: v.G, B: v.B, O: o}, // left top front
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 3, R: v.R, G: v.G, B: v.B, O: o}, // left top back
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 3, R: v.R, G: v.G, B: v.B, O: o}, // right top front
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 3, R: v.R, G: v.G, B: v.B, O: o}, // right top front
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 3, R: v.R, G: v.G, B: v.B, O: o}, // left top back
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, N: 3, R: v.R, G: v.G, B: v.B, O: o}) // right top back
	}

	// Bottom faces (Y-) N=4
	if yn {
		o := brightness(x, y-1, z)
		data = append(data,
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, N: 4, R: v.R, G: v.G, B: v.B, O: o}, // left
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 4, R: v.R, G: v.G, B: v.B, O: o}, // right
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 4, R: v.R, G: v.G, B: v.B, O: o}, //
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 4, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 4, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 4, R: v.R, G: v.G, B: v.B, O: o})
	}

	// Front faces (Z+) N=5
	if zp {
		o := brightness(x, y, z+1)
		data = append(data,
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B, O: o})
	}

	// Back faces (Z-) N=6
	if zn {
		o := brightness(x, y, z-1)
		data = append(data,
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B, O: o},
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B, O: o})
	}

	return data
}
