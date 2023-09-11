package terrain

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/math/ivec2"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/noise"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	*mesh.Dynamic[Vertex, uint16]
	Tile *Tile
}

func NewMesh(tile *Tile) *Mesh {
	mat := &material.Def{
		Pass:         material.Deferred,
		Shader:       "deferred/terrain",
		VertexFormat: Vertex{},
		DepthTest:    true,
		DepthWrite:   true,
		Primitive:    vertex.Triangles,
		CullMode:     vertex.CullBack,
		Transparent:  true, // does not cast shadows
	}
	msh := mesh.NewDynamic("Terrain", mat, FlatTileGenerator(tile, 0))
	msh.SetTexture("pattern", texture.Checker)
	msh.SetTexture("diffuse0", noise.NewWhiteNoise(64, 64))
	msh.SetTexture("diffuse1", noise.NewWhiteNoise(256, 256))

	tile.Changed.Subscribe(func(t *Tile) {
		msh.Refresh()
	})

	return &Mesh{
		Dynamic: msh,
		Tile:    tile,
	}
}

var normSamples = []ivec2.T{
	{X: 0, Y: 1},
	{X: 1, Y: 1},
	{X: 1, Y: 0},
	{X: 1, Y: -1},
	{X: 0, Y: -1},
	{X: -1, Y: -1},
	{X: -1, Y: 0},
	{X: -1, Y: 1},
	{X: 0, Y: 1},
}

func SmoothTileGenerator(tile *Tile) mesh.Generator[Vertex, uint16] {
	if tile.Size > 100 {
		panic("tile size cant be greater than 100x100")
	}
	return func() mesh.Data[Vertex, uint16] {
		side := tile.Size + 1

		getPoint := func(x, z int) (Point, bool) {
			tx, tz := (x+tile.Size)%tile.Size, (z+tile.Size)%tile.Size
			ox, oz := (x+tile.Size)/tile.Size-1, (z+tile.Size)/tile.Size-1
			t := tile.Map.Tile(tile.Position.X+ox, tile.Position.Y+oz, false)
			if t == nil {
				return Point{}, false
			}
			return t.Point(tx, tz), true
		}

		getVertex := func(x, z int) Vertex {
			root, _ := getPoint(x, z)
			origin := vec3.New(float32(x), root.Height, float32(z))

			norm := vec3.Zero
			samples := 0
			for i := 0; i < 8; i++ {
				ao := normSamples[i]
				ap, ok := getPoint(x+ao.X, z+ao.Y)
				if !ok {
					continue
				}
				a := vec3.New(float32(x+ao.X), ap.Height, float32(z+ao.Y)).Sub(origin)

				bo := normSamples[i+1]
				bp, ok := getPoint(x+bo.X, z+bo.Y)
				if !ok {
					continue
				}
				b := vec3.New(float32(x+bo.X), bp.Height, float32(z+bo.Y)).Sub(origin)

				norm = norm.Add(vec3.Cross(a, b).Normalized())
				samples++
			}

			norm = norm.Scaled(float32(1) / float32(samples))
			return Vertex{
				P: vec3.New(float32(x), root.Height, float32(z)),
				T: vec2.New(float32(x)/float32(tile.Size), 1-float32(z)/float32(tile.Size)),
				N: norm,
				W: vec4.New(root.Weights[0], root.Weights[1], root.Weights[2], root.Weights[3]),
			}
		}

		// generate vertices
		vertices := make([]Vertex, 0, side*side)
		indices := make([]uint16, 0, tile.Size*tile.Size*6)
		for z := 0; z < side; z++ {
			for x := 0; x < side; x++ {
				v := getVertex(x, z)
				vertices = append(vertices, v)
			}
		}

		// generate face indices
		idx := func(x, z int) uint16 {
			return uint16(z*side + x)
		}
		for z := 0; z < tile.Size; z++ {
			for x := 0; x < tile.Size; x++ {
				v00 := idx(x, z)
				v01 := idx(x, z+1)
				v10 := idx(x+1, z)
				v11 := idx(x+1, z+1)

				ex, ez := x%2 == 0, z%2 == 0
				if ex == ez {
					indices = append(indices, v00, v11, v10)
					indices = append(indices, v00, v01, v11)
				} else {
					indices = append(indices, v00, v01, v10)
					indices = append(indices, v01, v11, v10)
				}
			}
		}

		return mesh.Data[Vertex, uint16]{
			Vertices: vertices,
			Indices:  indices,
		}
	}
}

func FlatTileGenerator(tile *Tile, levelOfDetail int) mesh.Generator[Vertex, uint16] {
	if tile.Size < 1 {
		panic("tile size must be greater than 0")
	}
	if !IsPowerOfTwo(tile.Size) {
		panic("tile size must be a power of 2")
	}

	step := 1 << levelOfDetail
	if step >= tile.Size {
		panic("level of detail is too high for the tile size")
	}

	return func() mesh.Data[Vertex, uint16] {

		// generate vertices
		vertices := make([]Vertex, 0, 6*tile.Size*tile.Size)
		vertex := func(x, z int) Vertex {
			root := tile.Point(x, z)
			return Vertex{
				P: vec3.New(float32(x), root.Height, float32(z)),
				T: vec2.New(float32(x)/float32(tile.Size), 1-float32(z)/float32(tile.Size)),
				W: vec4.New(root.Weights[0], root.Weights[1], root.Weights[2], root.Weights[3]),
			}
		}

		steps := tile.Size / step
		for z := 0; z < steps; z++ {
			for x := 0; x < steps; x++ {
				ex, ez := x%2 == 0, z%2 == 0
				sx, sz := step*x, step*z
				if ex == ez {
					v1_00 := vertex(sx, sz)
					v1_10 := vertex(sx+step, sz)
					v1_11 := vertex(sx+step, sz+step)

					v2_00 := v1_00
					v2_01 := vertex(sx, sz+step)
					v2_11 := v1_11

					n1 := vec3.Cross(v1_11.P.Sub(v1_00.P), v1_10.P.Sub(v1_00.P)).Normalized()
					v1_00.N, v1_10.N, v1_11.N = n1, n1, n1

					n2 := vec3.Cross(v2_01.P.Sub(v2_00.P), v2_11.P.Sub(v2_00.P)).Normalized()
					v2_00.N, v2_01.N, v2_11.N = n2, n2, n2

					vertices = append(vertices, v1_00, v1_11, v1_10)
					vertices = append(vertices, v2_00, v2_01, v2_11)
				} else {
					v1_00 := vertex(sx, sz)
					v1_01 := vertex(sx, sz+step)
					v1_10 := vertex(sx+step, sz)

					v2_01 := v1_01
					v2_11 := vertex(sx+step, sz+step)
					v2_10 := v1_10

					n1 := vec3.Cross(v1_01.P.Sub(v1_00.P), v1_10.P.Sub(v1_00.P)).Normalized()
					v1_00.N, v1_10.N, v1_01.N = n1, n1, n1

					n2 := vec3.Cross(v2_11.P.Sub(v2_01.P), v2_10.P.Sub(v2_01.P)).Normalized()
					v2_01.N, v2_11.N, v2_10.N = n2, n2, n2

					vertices = append(vertices, v1_00, v1_01, v1_10)
					vertices = append(vertices, v2_01, v2_11, v2_10)
				}
			}
		}

		return mesh.Data[Vertex, uint16]{
			Vertices: vertices,
			Indices:  nil,
		}
	}
}

func IsPowerOfTwo(x int) bool {
	return (x & (x - 1)) == 0
}
