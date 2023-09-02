package terrain

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/math/ivec2"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/math/vec4"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/noise"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	*mesh.Dynamic[Vertex, uint16]
	Tile *Tile
}

func NewMesh(tile *Tile) *Mesh {
	msh := mesh.NewDynamic(
		"Terrain",
		&material.Def{
			Pass:         material.Deferred,
			Shader:       "deferred/terrain",
			VertexFormat: Vertex{},
			DepthTest:    true,
			DepthWrite:   true,
			Primitive:    vertex.Triangles,
			CullMode:     vertex.CullBack,
		},
		TileVertexGenerator(tile),
	)
	// msh.SetTexture(texture.Diffuse, texture.Checker)
	msh.SetTexture("pattern", noise.NewWhiteNoise(64, 64))
	msh.SetTexture("diffuse0", noise.NewWhiteNoise(256, 256))
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

func TileVertexGenerator(tile *Tile) mesh.Generator[Vertex, uint16] {
	if tile.Size > 100 {
		panic("tile size cant be greater than 100x100")
	}
	return func() mesh.Data[Vertex, uint16] {
		side := tile.Size + 1

		getPoint := func(x, z int) (Point, bool) {
			tx, tz := (x+tile.Size)%tile.Size, (z+tile.Size)%tile.Size
			ox, oz := (x+tile.Size)/tile.Size-1, (z+tile.Size)/tile.Size-1
			t := tile.Map.GetTile(tile.Position.X+ox, tile.Position.Y+oz, false)
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
				W: vec4.New(1, 0, 0, 0),
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
				indices = append(indices, v00, v11, v10)
				indices = append(indices, v00, v01, v11)
			}
		}

		return mesh.Data[Vertex, uint16]{
			Vertices: vertices,
			Indices:  indices,
		}
	}
}
