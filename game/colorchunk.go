package game

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/render"
)

/* Chunks are smallest renderable units of voxel geometry */
type ColorChunk struct {
	*engine.ComponentBase

	Size int
	Data ColorVoxels

	vao  *render.VertexArray
	vbo  *render.VertexBuffer
	mesh *engine.Mesh
}

type OcclusionSpace interface {
	Free(x, y, z int) bool
}

type ColorVoxel struct {
	R, G, B byte
}

func NewColorVoxel(color render.Color) *ColorVoxel {
	return &ColorVoxel{
		R: byte(255 * color.R),
		G: byte(255 * color.G),
		B: byte(255 * color.B),
	}
}

type ColorVoxelVertex struct {
	X, Y, Z byte // position
	N       byte // normal index
	R, G, B byte // color
	O       byte // occlusion
}

// todo: dont use pointers
type ColorVoxels []*ColorVoxel

// vertex array type
type ColorVoxelVertices []ColorVoxelVertex

func (buffer ColorVoxelVertices) Elements() int { return len(buffer) }
func (buffer ColorVoxelVertices) Size() int     { return 8 }

func NewColorChunk(parentObject *engine.Object, size int) *ColorChunk {
	chk := &ColorChunk{
		Size: size,
		Data: make(ColorVoxels, size*size*size),

		vao: render.CreateVertexArray(),
		vbo: render.CreateVertexBuffer(),
	}

	chk.ComponentBase = engine.NewComponent(parentObject, chk)
	return chk
}

/* Clears all voxel data in this chunk */
func (chk *ColorChunk) Clear() {
	for i := 0; i < len(chk.Data); i++ {
		chk.Data[i] = nil
	}
}

/* Returns the slice offset for a given set of coordinates, as
   well as a bool indicating whether the position is within bounds.
   If the point is out of bounds, zero is returned */
func (chk *ColorChunk) offset(x, y, z int) (int, bool) {
	s := chk.Size
	if x < 0 || x >= s || y < 0 || y >= s || z < 0 || z >= s {
		return 0, false
	}
	s2 := s * s
	pos := z*s2 + y*s + x
	return pos, true
}

/* Returns a pointer to the voxel defintion at the given position.
   If the space is empty, nil is returned */
func (chk *ColorChunk) At(x, y, z int) *ColorVoxel {
	pos, ok := chk.offset(x, y, z)
	if !ok {
		return nil
	}
	return chk.Data[pos]
}

/* Sets a voxel. If it is outside bounds, nothing happens */
func (chk *ColorChunk) Set(x, y, z int, voxel *ColorVoxel) {
	pos, ok := chk.offset(x, y, z)
	if !ok {
		return
	}
	chk.Data[pos] = voxel
}

func (chk *ColorChunk) Free(x, y, z int) bool {
	v, ok := chk.offset(x, y, z)
	if !ok {
		return true
	}
	return chk.Data[v] == nil
}

func (chk *ColorChunk) Update(dt float32) {
}

func (chk *ColorChunk) Draw(args render.DrawArgs) {
	if args.Pass == "geometry" {
		chk.vao.Draw()
	}
}

type OcclusionData struct {
	data   []float32
	length int
	Size   int
}

func NewOcclusionData(size int) *OcclusionData {
	l := size * size * size
	return &OcclusionData{
		Size:   size,
		length: l,
		data:   make([]float32, l),
	}
}

func (o *OcclusionData) Get(x, y, z byte) byte {
	if x < 0 || y < 0 || z < 0 || int(x) >= o.Size || int(y) >= o.Size || int(z) >= o.Size {
		return 255
	}
	offset := int(z)*o.Size*o.Size + int(y)*o.Size + int(x)
	return byte(256 * o.data[offset])
}
func (o *OcclusionData) Set(x, y, z int, value float32) {
	if x < 0 || y < 0 || z < 0 || int(x) >= o.Size || int(y) >= o.Size || int(z) >= o.Size {
		return
	}
	offset := int(z)*o.Size*o.Size + int(y)*o.Size + int(x)
	o.data[offset] = value
}

/* Recomputes the chunk mesh and returns a pointer to it. */
func (chk *ColorChunk) Compute() {
	s := chk.Size
	data := make(ColorVoxelVertices, 0, 64)

	occlusion := NewOcclusionData(s)
	f := func(f bool) float32 {
		if f {
			return 1.0 / 6
		}
		return 0
	}

	/* occlusion pass */
	for z := 0; z < s; z++ {
		for y := 0; y < s; y++ {
			for x := 0; x < s; x++ {
				v := chk.At(x, y, z)
				if v != nil {
					/* Not empty space */
					continue
				}

				/* Simple optimization - dont draw hidden faces */
				xp := chk.At(x+1, y, z) == nil
				xn := chk.At(x-1, y, z) == nil
				yp := chk.At(x, y+1, z) == nil
				yn := chk.At(x, y-1, z) == nil
				zp := chk.At(x, y, z+1) == nil
				zn := chk.At(x, y, z-1) == nil

				o := f(xp) + f(xn) + f(yp) + f(yn) + f(zp) + f(zn)
				occlusion.Set(x, y, z, o)
			}
		}
	}

	/* geometry pass */
	for z := 0; z < s; z++ {
		for y := 0; y < s; y++ {
			for x := 0; x < s; x++ {
				v := chk.At(x, y, z)
				if v == nil {
					/* Empty space */
					continue
				}

				/* Simple optimization - dont draw hidden faces */
				xp := chk.At(x+1, y, z) == nil
				xn := chk.At(x-1, y, z) == nil
				yp := chk.At(x, y+1, z) == nil
				yn := chk.At(x, y-1, z) == nil
				zp := chk.At(x, y, z+1) == nil
				zn := chk.At(x, y, z-1) == nil

				/* Compute & append vertex data */
				vertices := v.Compute(occlusion, byte(x), byte(y), byte(z), xp, xn, yp, yn, zp, zn)
				data = append(data, vertices...)
			}
		}
	}

	/* Buffer to GPU */
	chk.vao.Length = int32(len(data))
	if chk.vao.Length > 0 {
		chk.vao.Bind()
		chk.vbo.Buffer(data)
	}
}

func (v ColorVoxel) Compute(occlusion *OcclusionData, x, y, z byte, xp, xn, yp, yn, zp, zn bool) ColorVoxelVertices {
	data := make(ColorVoxelVertices, 0, 36)

	// Right (X+) N=1
	if xp {
		o := occlusion.Get(x+1, y, z)
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
		o := occlusion.Get(x-1, y, z)
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
		o := occlusion.Get(x, y+1, z)
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
		o := occlusion.Get(x, y-1, z)
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
		o := occlusion.Get(x, y, z+1)
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
		o := occlusion.Get(x, y, z-1)
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
