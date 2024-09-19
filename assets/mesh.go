package assets

import (
	"github.com/johanhenriksson/goworld/assets/fs"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Mesh interface {
	Asset

	// LoadMesh is called by mesh caches and loaders, and should return the mesh data.
	LoadMesh(assets fs.Filesystem) vertex.Mesh
}
