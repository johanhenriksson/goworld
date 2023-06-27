package mesh

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Objects that implement the mesh.UpdateHandler interface will receive
// a callback if a sibling Mesh updates its mesh data.
// This allows other components to react to changing meshes.
type UpdateHandler interface {
	object.T
	OnMeshUpdate(vertex.Mesh)
}
