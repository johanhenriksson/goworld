package game

import (
	"github.com/johanhenriksson/goworld/render"
)

// EmptyVoxel is an empty color voxel
var EmptyVoxel = Voxel{}

// Voxels is a collection of voxels
type Voxels []Voxel

// Voxel holds color information for a single colored voxel
type Voxel struct {
	R, G, B byte
}

// NewVoxel creates a new Color Voxel from a given color
func NewVoxel(color render.Color) Voxel {
	return Voxel{
		R: byte(255 * color.R),
		G: byte(255 * color.G),
		B: byte(255 * color.B),
	}
}

var LightSamples = [][][]int{
	{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}},         // N=0
	{{0, -1, -1}, {0, -1, 0}, {0, 0, -1}, {0, 0, 0}},     // X+
	{{-1, -1, -1}, {-1, -1, 0}, {-1, 0, -1}, {-1, 0, 0}}, // X-
	{{-1, 0, -1}, {-1, 0, 0}, {0, 0, 0}, {0, 0, -1}},     // Y+
	{{-1, -1, -1}, {-1, -1, 0}, {0, -1, 0}, {0, -1, -1}}, // Y-
	{{-1, -1, 0}, {0, -1, 0}, {-1, 0, 0}, {0, 0, 0}},     // Z+
	{{-1, -1, -1}, {0, -1, -1}, {-1, 0, -1}, {0, 0, -1}}, // Z-
}

// Compute the verticies of this Color Voxel
func (v Voxel) Compute(light *LightVolume, x, y, z byte, xp, xn, yp, yn, zp, zn bool) VoxelVertices {
	data := make(VoxelVertices, 0, 36)

	sampleOcclusion := func(v *VoxelVertex) {
		sample := func(v *VoxelVertex, i int) float32 {
			of := LightSamples[v.N][i]
			l := light.Get(int(v.X)+of[0], int(v.Y)+of[1], int(v.Z)+of[2])
			if l != nil && !l.Blocked {
				return l.V
			}
			return 0
		}
		o := float32(0)
		o += sample(v, 0)
		o += sample(v, 1)
		o += sample(v, 2)
		o += sample(v, 3)
		v.O = byte(255 * (1 - o/float32(4)))
	}

	// Right (X+) N=1
	if xp {
		data = append(data,
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 1, R: v.R, G: v.G, B: v.B}, // front bottom right
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 1, R: v.R, G: v.G, B: v.B}, // back bottom right
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 1, R: v.R, G: v.G, B: v.B}, // back top right
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 1, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 1, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, N: 1, R: v.R, G: v.G, B: v.B})
	}

	// Left faces (X-) N=2
	if xn {
		data = append(data,
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 2, R: v.R, G: v.G, B: v.B}, // bottom left back
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 2, R: v.R, G: v.G, B: v.B}, // top left front
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, N: 2, R: v.R, G: v.G, B: v.B}, // bottom left front
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 2, R: v.R, G: v.G, B: v.B}, // bottom left back
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 2, R: v.R, G: v.G, B: v.B}, // top left back
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 2, R: v.R, G: v.G, B: v.B}) // top left front
	}

	// Top faces (Y+) N=3
	if yp {
		data = append(data,
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 3, R: v.R, G: v.G, B: v.B}, // left top front
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 3, R: v.R, G: v.G, B: v.B}, // left top back
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 3, R: v.R, G: v.G, B: v.B}, // right top front
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 3, R: v.R, G: v.G, B: v.B}, // right top front
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 3, R: v.R, G: v.G, B: v.B}, // left top back
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, N: 3, R: v.R, G: v.G, B: v.B}) // right top back
	}

	// Bottom faces (Y-) N=4
	if yn {
		data = append(data,
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, N: 4, R: v.R, G: v.G, B: v.B}, // left
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 4, R: v.R, G: v.G, B: v.B}, // right
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 4, R: v.R, G: v.G, B: v.B}, //
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 4, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 4, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 4, R: v.R, G: v.G, B: v.B})
	}

	// Front faces (Z+) N=5
	if zp {
		data = append(data,
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 1, N: 5, R: v.R, G: v.G, B: v.B})
	}

	// Back faces (Z-) N=6
	if zn {
		data = append(data,
			VoxelVertex{X: x + 0, Y: y + 0, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 1, Y: y + 0, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 0, Y: y + 1, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B},
			VoxelVertex{X: x + 1, Y: y + 1, Z: z + 0, N: 6, R: v.R, G: v.G, B: v.B})
	}

	for i := range data {
		sampleOcclusion(&data[i])
	}

	return data
}
