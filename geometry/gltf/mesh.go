package gltf

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type gltfMesh struct {
	key       string
	elements  int
	primitive vertex.Primitive
	pointers  vertex.Pointers
	indices   []byte
	vertices  []byte
	indexsize int
}

var _ vertex.Mesh = &gltfMesh{}

func (m *gltfMesh) Key() string      { return m.key }
func (m *gltfMesh) Version() int     { return 1 }
func (m *gltfMesh) IndexCount() int  { return m.elements }
func (m *gltfMesh) IndexSize() int   { return m.indexsize }
func (m *gltfMesh) IndexData() any   { return m.indices }
func (m *gltfMesh) VertexCount() int { return len(m.vertices) / m.VertexSize() }
func (m *gltfMesh) VertexSize() int  { return m.pointers.Stride() }
func (m *gltfMesh) VertexData() any  { return m.vertices }

func (m *gltfMesh) Positions(func(vec3.T)) {
	panic("not implemented")
}

func (m *gltfMesh) Min() vec3.T { panic("not implemented") }
func (m *gltfMesh) Max() vec3.T { panic("not implemented") }

func (m *gltfMesh) Primitive() vertex.Primitive { return m.primitive }
func (m *gltfMesh) Pointers() vertex.Pointers   { return m.pointers }
