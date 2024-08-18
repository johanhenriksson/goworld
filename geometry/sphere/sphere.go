package sphere

import (
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh struct {
	*mesh.Static
	Subdivisions object.Property[int]

	data vertex.MutableMesh[vertex.T, uint16]
}

func New(pool object.Pool, mat *material.Def) *Mesh {
	m := object.NewComponent(pool, &Mesh{
		Static:       mesh.New(pool, mat),
		Subdivisions: object.NewProperty(3),
	})
	m.SetTexture(texture.Diffuse, texture.Checker)
	m.data = vertex.NewTriangles[vertex.T, uint16](object.Key("sphere", m), nil, nil)
	m.Subdivisions.OnChange.Subscribe(func(int) { m.refresh() })
	m.refresh()
	return m
}

func (m *Mesh) refresh() {
	tris := icosphere(m.Subdivisions.Get())

	vertices := []vertex.T{}
	for _, tri := range tris {
		vertices = append(vertices, vertex.T{
			P: tri.A,
			N: tri.A,
			T: vec2.New(0, 0),
		})
		vertices = append(vertices, vertex.T{
			P: tri.B,
			N: tri.B,
			T: vec2.New(0, 0),
		})
		vertices = append(vertices, vertex.T{
			P: tri.C,
			N: tri.C,
			T: vec2.New(0, 0),
		})
	}

	m.data.Update(vertices, nil)
	m.VertexData.Set(m.data)
}

func icosphere(subdivisions int) []vertex.Triangle {
	const X = float32(0.525731112119133606)
	const Z = float32(0.850650808352039932)

	vertices := []vec3.T{
		vec3.New(-X, 0, Z),
		vec3.New(X, 0, Z),
		vec3.New(-X, 0, -Z),
		vec3.New(X, 0, -Z),
		vec3.New(0, Z, X),
		vec3.New(0, Z, -X),
		vec3.New(0, -Z, X),
		vec3.New(0, -Z, -X),
		vec3.New(Z, X, 0),
		vec3.New(-Z, X, 0),
		vec3.New(Z, -X, 0),
		vec3.New(-Z, -X, 0),
	}

	faces := []vertex.Triangle{
		{A: vertices[1], B: vertices[4], C: vertices[0]},
		{A: vertices[4], B: vertices[9], C: vertices[0]},
		{A: vertices[4], B: vertices[5], C: vertices[9]},
		{A: vertices[8], B: vertices[5], C: vertices[4]},
		{A: vertices[1], B: vertices[8], C: vertices[4]},
		{A: vertices[1], B: vertices[10], C: vertices[8]},
		{A: vertices[10], B: vertices[3], C: vertices[8]},
		{A: vertices[8], B: vertices[3], C: vertices[5]},
		{A: vertices[3], B: vertices[2], C: vertices[5]},
		{A: vertices[3], B: vertices[7], C: vertices[2]},
		{A: vertices[3], B: vertices[10], C: vertices[7]},
		{A: vertices[10], B: vertices[6], C: vertices[7]},
		{A: vertices[6], B: vertices[11], C: vertices[7]},
		{A: vertices[6], B: vertices[0], C: vertices[11]},
		{A: vertices[6], B: vertices[1], C: vertices[0]},
		{A: vertices[10], B: vertices[1], C: vertices[6]},
		{A: vertices[11], B: vertices[0], C: vertices[9]},
		{A: vertices[2], B: vertices[11], C: vertices[9]},
		{A: vertices[5], B: vertices[2], C: vertices[9]},
		{A: vertices[11], B: vertices[2], C: vertices[7]},
	}

	for r := subdivisions; r > 0; r-- {
		result := make([]vertex.Triangle, 0, 4*len(faces))
		for _, tri := range faces {
			v1 := vec3.Mid(tri.A, tri.B).Normalized()
			v2 := vec3.Mid(tri.B, tri.C).Normalized()
			v3 := vec3.Mid(tri.C, tri.A).Normalized()
			result = append(result, vertex.Triangle{A: tri.A, B: v1, C: v3})
			result = append(result, vertex.Triangle{A: tri.B, B: v2, C: v1})
			result = append(result, vertex.Triangle{A: tri.C, B: v3, C: v2})
			result = append(result, vertex.Triangle{A: v1, B: v2, C: v3})
		}
		faces = result
	}

	return faces
}
