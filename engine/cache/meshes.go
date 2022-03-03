package cache

import (
	"log"

	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Meshes interface {
	// Fetches the given mesh & material combination from the cache
	// If it does not exist, it will be inserted.
	Fetch(vertex.Mesh, material.T) GpuMesh

	// Evict the first cached mesh that is older than max age
	// Returns true if a mesh is evicted
	Evict() bool

	// Tick increments the age of all items in the cache
	Tick()
}

type meshes struct {
	maxAge  int
	cache   map[string]*entry
	backend MeshBackend
}

type entry struct {
	age     int
	version int
	mesh    GpuMesh
}

func NewMeshes() Meshes {
	return &meshes{
		maxAge:  1000,
		cache:   make(map[string]*entry),
		backend: &glmeshes{},
	}
}

func (m *meshes) Fetch(mesh vertex.Mesh, mat material.T) GpuMesh {
	// if the mesh is nil, just return a no-op mesh
	if mesh == nil {
		return nilmesh{}
	}

	line, hit := m.cache[mesh.Id()]

	// not in cache - instantiate a buffered mesh
	if !hit {
		log.Println("buffering new mesh", mesh.Id(), "version", mesh.Version())
		vao := m.backend.Instantiate(mesh, mat)
		line = &entry{
			mesh: vao,
		}
		m.cache[mesh.Id()] = line
	}

	// version has changed, update the mesh
	if line.version != mesh.Version() {
		// we might want to queue this operation and run it at a more appropriate time
		log.Println("updating existing mesh", mesh.Id(), "to version", mesh.Version())
		m.backend.Update(line.mesh, mesh)
		line.version = mesh.Version()
	}

	// reset age
	line.age = 0

	return line.mesh
}

func (m *meshes) Tick() {
	// increment the age of every item in the cache
	for _, entry := range m.cache {
		entry.age++
	}
}

func (m *meshes) Evict() bool {
	for id, entry := range m.cache {
		// skip any meshes that have been recently used
		if entry.age < m.maxAge {
			continue
		}

		// deallocate gpu memory
		log.Println("deallocating", id, "from gpu")
		m.backend.Delete(m.cache[id].mesh)

		// remove cache line
		delete(m.cache, id)
		return true
	}
	return false
}
