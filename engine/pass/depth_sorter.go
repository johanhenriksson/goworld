package pass

import (
	"sort"

	"github.com/johanhenriksson/goworld/core/mesh"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

func DepthSortGroups(cache MatCache, args render.Args, meshes []mesh.Mesh) DrawGroups {
	// perform back-to-front depth sorting
	// we use the closest point on the meshes bounding sphere as a heuristic
	sort.SliceStable(meshes, func(i, j int) bool {
		// return true if meshes[i] is further away than meshes[j]
		first, second := meshes[i].BoundingSphere(), meshes[j].BoundingSphere()

		di := vec3.Distance(args.Position, first.Center) - first.Radius
		dj := vec3.Distance(args.Position, second.Center) - second.Radius
		return di > dj
	})

	// sort meshes by material
	return OrderedGroups(cache, args.Context.Index, meshes)
}
