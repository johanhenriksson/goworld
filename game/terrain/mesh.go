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
	msh.SetTexture("heightmap", texture.PathRef("textures/uv_checker.png"))
	return &Mesh{
		Dynamic: msh,
		Tile:    tile,
	}
}

func TileVertexGenerator(tile *Tile) mesh.Generator[vertex.T, uint16] {
	return func() mesh.Data[vertex.T, uint16] {
		side := tile.Size + 1
		vertices := make([]vertex.T, 0, side*side)
		for y := 0; y < side; y++ {
			for x := 0; x < side; x++ {
				vertices = append(vertices, vertex.T{
					P: vec3.NewI(x, 0, y),
					T: vec2.New(float32(x)/float32(tile.Size), 1-float32(y)/float32(tile.Size)),
				})
			}
		}

		indices := make([]uint16, 0, tile.Size*tile.Size*6)
		for y := 0; y < tile.Size; y++ {
			for x := 0; x < tile.Size; x++ {
				// t1
				indices = append(indices, uint16(side*(y)+x+1))
				indices = append(indices, uint16(side*(y)+x))
				indices = append(indices, uint16(side*(y+1)+x))

				// t2
				indices = append(indices, uint16(side*(y+1)+x))
				indices = append(indices, uint16(side*(y+1)+x+1))
				indices = append(indices, uint16(side*(y)+x+1))
			}
		}

		return mesh.Data[vertex.T, uint16]{
			Vertices: vertices,
			Indices:  indices,
		}
	}
}
