package cache

import (
	"log"

	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type BufferedMesh interface {
	Delete()
	Draw() error
}

type nilmesh struct{}

func (n nilmesh) Draw() error { return nil }
func (n nilmesh) Delete()     {}

type Meshes interface {
	Fetch(vertex.Mesh, material.T) BufferedMesh
	Evict()
}

type entry struct {
	age  int
	mesh BufferedMesh
}

type meshes struct {
	cache  map[string]*entry
	evicts []string
}

func NewMeshes() Meshes {
	return &meshes{
		cache:  make(map[string]*entry),
		evicts: make([]string, 0, 64),
	}
}

func (m *meshes) Fetch(mesh vertex.Mesh, mat material.T) BufferedMesh {
	if mesh == nil {
		return nilmesh{}
	}

	line, hit := m.cache[mesh.Id()]
	if hit {
		line.age--
		return line.mesh
	}

	// this is where we create a vao!
	vao := gl_vertex_array.New(mesh.Primitive())

	ptrs := mesh.Pointers()
	ptrs.Bind(mat)
	vao.SetPointers(ptrs)
	vao.SetIndexSize(mesh.IndexSize())
	vao.SetElements(mesh.Elements())
	vao.Buffer("vertex", mesh.VertexData())
	vao.Buffer("index", mesh.IndexData())

	log.Println("buffered mesh", mesh.Id(), "to gpu")
	m.cache[mesh.Id()] = &entry{
		mesh: vao,
	}
	return vao
}

func (m *meshes) Evict() {
	for id, entry := range m.cache {
		entry.age++
		if entry.age > 1000 {
			m.evicts = append(m.evicts, id)
		}
	}

	for _, id := range m.evicts {
		log.Println("deallocating", id, "from gpu")
		m.cache[id].mesh.Delete()
		delete(m.cache, id)
	}

	m.evicts = m.evicts[:0]
}
