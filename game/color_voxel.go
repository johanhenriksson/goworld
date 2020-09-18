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

	sampleOcclusion := func(v *ColorVoxelVertex) {
		o, n := float32(0), 0
		sample := func(x, y, z byte) {
			l := light.Get(int(x), int(y), int(z))
			if l != nil && !l.Blocked {
				n++
				o += l.V
			}
		}
		sample(v.X, v.Y, v.Z)
		sample(v.X-1, v.Y, v.Z)
		sample(v.X, v.Y-1, v.Z)
		sample(v.X-1, v.Y-1, v.Z)
		sample(v.X, v.Y, v.Z-1)
		sample(v.X-1, v.Y, v.Z-1)
		sample(v.X, v.Y-1, v.Z-1)
		sample(v.X-1, v.Y-1, v.Z-1)
		v.O = byte(255 * (1 - o/float32(n)))
	}

	// Right (X+) N=1
	if xp {
		data = append(data,
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 1, R: v.R, G: v.G, B: v.B}, // front bottom right
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 1, R: v.R, G: v.G, B: v.B}, // back bottom right
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 1, R: v.R, G: v.G, B: v.B}, // back top right
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 1, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 1, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, N: 1, R: v.R, G: v.G, B: v.B})
	}

	// Left faces (X-) N=2
	if xn {
		data = append(data,
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 2, R: v.R, G: v.G, B: v.B}, // bottom left back
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 2, R: v.R, G: v.G, B: v.B}, // top left front
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, N: 2, R: v.R, G: v.G, B: v.B}, // bottom left front
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 2, R: v.R, G: v.G, B: v.B}, // bottom left back
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 2, R: v.R, G: v.G, B: v.B}, // top left back
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 2, R: v.R, G: v.G, B: v.B}) // top left front
	}

	// Top faces (Y+) N=3
	if yp {
		data = append(data,
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 3, R: v.R, G: v.G, B: v.B}, // left top front
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 3, R: v.R, G: v.G, B: v.B}, // left top back
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 3, R: v.R, G: v.G, B: v.B}, // right top front
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 3, R: v.R, G: v.G, B: v.B}, // right top front
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 3, R: v.R, G: v.G, B: v.B}, // left top back
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, N: 3, R: v.R, G: v.G, B: v.B}) // right top back
	}

	// Bottom faces (Y-) N=4
	if yn {
		data = append(data,
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, N: 4, R: v.R, G: v.G, B: v.B}, // left
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 4, R: v.R, G: v.G, B: v.B}, // right
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 4, R: v.R, G: v.G, B: v.B}, //
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 4, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 4, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 4, R: v.R, G: v.G, B: v.B})
	}

	// Front faces (Z+) N=5
	if zp {
		data = append(data,
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B})
	}

	// Back faces (Z-) N=6
	if zn {
		data = append(data,
			ColorVoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B},
			ColorVoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B})
	}

	for i := range data {
		sampleOcclusion(&data[i])
	}

	return data
}
