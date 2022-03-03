package cache

import (
	"log"

	"github.com/johanhenriksson/goworld/render/backend/gl/gl_vertex_array"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type BufferedMesh interface {
	Draw() error
}

type Meshes interface {
	Fetch(vertex.Mesh, material.T) BufferedMesh
	Evict()
}

type entry struct {
	age     int
	version int
	mesh    vertex.Array
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
	if !hit {
		log.Println("buffering new mesh", mesh.Id(), "version", mesh.Version())
		return m.instantiate(mesh, mat)
	}

	// decrement age (frame since last use)
	line.age--

	// if we have the appropriate version, just return it
	if line.version == mesh.Version() {
		return line.mesh
	}

	// version has changed, update the mesh
	log.Println("updating existing mesh", mesh.Id(), "to version", mesh.Version())
	return m.update(line, mesh, mat)
}

func (m *meshes) instantiate(mesh vertex.Mesh, mat material.T) BufferedMesh {
	// most of this is opengl-specific and could be extracted to make the cache more general

	vao := gl_vertex_array.New(mesh.Primitive())

	line := &entry{
		mesh: vao,
	}
	m.cache[mesh.Id()] = line
	return m.update(line, mesh, mat)
}

func (m *meshes) update(line *entry, mesh vertex.Mesh, mat material.T) BufferedMesh {
	// most of this is opengl-specific and could be extracted to make the cache more general

	ptrs := mesh.Pointers()
	ptrs.Bind(mat)

	line.mesh.SetPointers(ptrs)
	line.mesh.SetIndexSize(mesh.IndexSize())
	line.mesh.SetElements(mesh.Elements())
	line.mesh.Buffer("vertex", mesh.VertexData())
	line.mesh.Buffer("index", mesh.IndexData())
	line.version = mesh.Version()

	return line.mesh
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
