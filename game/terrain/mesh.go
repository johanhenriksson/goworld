package terrain

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	mesh.Dynamic[vertex.T, uint16]
	Tile *Tile
}

func NewMesh(tile *Tile) *Mesh {
	msh := mesh.NewDynamic("Terrain", mesh.Deferred, &material.Def{
		Shader:       "game/terrain",
		Subpass:      "geometry",
		VertexFormat: vertex.T{},
		DepthTest:    true,
		DepthWrite:   true,
		CullMode:     vertex.CullBack,
	}, TileVertexGenerator(tile))
	msh.SetTexture("heightmap", texture.PathRef("textures/heightmap.png"))
	msh.SetTexture("diffuse", texture.PathRef("textures/uv_checker.png"))
	return &Mesh{
		Dynamic: msh,
		Tile:    tile,
	}
}

func TileVertexGenerator(tile *Tile) mesh.Generator[vertex.T, uint16] {
	return func() mesh.Data[vertex.T, uint16] {
		side := tile.Size + 1

		getVertex := func(x, z int) vertex.T {
			return vertex.T{
				P: vec3.New(float32(x), tile.points[z][x].Height, float32(z)),
				T: vec2.New(float32(x)/float32(tile.Size), 1-float32(z)/float32(tile.Size)),
			}
		}

		// generate two faces for each square
		vertices := make([]vertex.T, 0, side*side*6)
		addTriangle := func(v0, v1, v2 vertex.T) {
			a := v1.P.Sub(v0.P)
			b := v2.P.Sub(v0.P)
			n := vec3.Cross(a, b).Normalized()
			v0.N = n
			v1.N = n
			v2.N = n
			vertices = append(vertices, v0, v1, v2)
		}

		for z := 0; z < tile.Size; z++ {
			for x := 0; x < tile.Size; x++ {
				v00 := getVertex(x, z)
				v01 := getVertex(x, z+1)
				v10 := getVertex(x+1, z)
				v11 := getVertex(x+1, z+1)

				addTriangle(v00, v11, v10)
				addTriangle(v00, v01, v11)
			}
		}

		indices := make([]uint16, 0, len(vertices))
		for i := range indices {
			indices[i] = uint16(i)
		}

		return mesh.Data[vertex.T, uint16]{
			Vertices: vertices,
			Indices:  indices,
		}
	}
}
