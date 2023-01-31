package collider

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/physics"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type T interface {
	object.T

	Intersect(ray *physics.Ray) (bool, vec3.T)
}
