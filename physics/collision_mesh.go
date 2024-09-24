package physics

import (
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Vertex format used by physics meshes
type Vertex struct {
	vec3.T `vtx:"position,float,3"`
}

func (v Vertex) Position() vec3.T {
	return v.T
}

func CollisionMesh(mesh vertex.Mesh) vertex.Mesh {
	// generate collision mesh
	// todo: use greedy face optimization

	indexMap := make(map[vec3.T]uint32, mesh.IndexCount())
	vertexdata := make([]Vertex, 0, mesh.VertexCount()/4)
	indexdata := make([]uint32, 0, mesh.IndexCount())
	for p := range mesh.Positions() {
		// check if the vertex position already has an index
		// todo: tolerance
		index, exists := indexMap[p]
		if !exists {
			// create a new index from the vertex
			index = uint32(len(vertexdata))
			vertexdata = append(vertexdata, Vertex{p})
			indexMap[p] = index
		}
		// store vertex index
		indexdata = append(indexdata, index)
	}

	return vertex.NewTriangles[Vertex, uint32](mesh.Key(), vertexdata, indexdata)
}
