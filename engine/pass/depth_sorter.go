package pass

import (
	"sort"

	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/engine/uniform"
	"github.com/johanhenriksson/goworld/math/vec3"
)

func DepthSortGroups(cache MaterialCache, frame int, cam uniform.Camera, meshes []mesh.Mesh) DrawGroups {
	eye := cam.Eye.XYZ()

	// perform back-to-front depth sorting
	// we use the closest point on the meshes bounding sphere as a heuristic
	sort.SliceStable(meshes, func(i, j int) bool {
		// return true if meshes[i] is further away than meshes[j]
		ri, rj := meshes[i].Mesh(), meshes[j].Mesh()

		// todo: LoadMesh is not allowed here
		mi, mj := ri.LoadMesh(assets.FS), rj.LoadMesh(assets.FS)
		if mi == nil || mj == nil {
			// doesnt matter which is in front
			return false
		}

		first := mi.Bounds(meshes[i].Transform().WorldPosition())
		second := mj.Bounds(meshes[j].Transform().WorldPosition())

		di := vec3.Distance(eye, first.Center) - first.Radius
		dj := vec3.Distance(eye, second.Center) - second.Radius
		return di > dj
	})

	// sort meshes by material
	return OrderedGroups(cache, frame, meshes)
}
