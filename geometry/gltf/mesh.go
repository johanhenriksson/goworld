package gltf

import (
	"github.com/johanhenriksson/goworld/render/vertex"
)

type gltfMesh struct {
	id        string
	elements  int
	primitive vertex.Primitive
	pointers  vertex.Pointers
	indices   []byte
	vertices  []byte
	indexsize int
}

func (m *gltfMesh) Id() string      { return m.id }
func (m *gltfMesh) Version() int    { return 1 }
func (m *gltfMesh) Elements() int   { return m.elements }
func (m *gltfMesh) IndexData() any  { return m.indices }
func (m *gltfMesh) VertexData() any { return m.vertices }
func (m *gltfMesh) IndexSize() int  { return m.indexsize }

func (m *gltfMesh) Primitive() vertex.Primitive { return m.primitive }
func (m *gltfMesh) Pointers() vertex.Pointers   { return m.pointers }
