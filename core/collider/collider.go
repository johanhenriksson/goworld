package collider

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/physics"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type T interface {
	object.Component

	Intersect(ray *physics.Ray) (bool, vec3.T)
}

func ClosestIntersection(colliders []T, ray *physics.Ray) (T, bool) {
	var closest T
	closestDist := float32(math.InfPos)
	for _, collider := range colliders {
		hit, point := collider.Intersect(ray)
		if hit {
			dist := vec3.Distance(point, ray.Origin)
			if dist < closestDist {
				closest = collider
				closestDist = dist
			}
		}
	}
	return closest, closest != nil
}
