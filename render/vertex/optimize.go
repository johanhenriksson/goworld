package vertex

import "github.com/johanhenriksson/goworld/math/vec3"

func CollisionMesh(mesh Mesh) Mesh {
	// generate collision mesh
	// todo: use greedy face optimization

	indexMap := make(map[vec3.T]uint32, mesh.IndexCount())
	vertexdata := make([]P, 0, mesh.VertexCount()/4)
	indexdata := make([]uint32, 0, mesh.IndexCount())
	for p := range mesh.Positions() {
		// check if the vertex position already has an index
		// todo: tolerance
		index, exists := indexMap[p]
		if !exists {
			// create a new index from the vertex
			index = uint32(len(vertexdata))
			vertexdata = append(vertexdata, P{p})
			indexMap[p] = index
		}
		// store vertex index
		indexdata = append(indexdata, index)
	}

	return NewTriangles[P, uint32](mesh.Key(), vertexdata, indexdata)
}
