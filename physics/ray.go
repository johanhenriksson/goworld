package physics

import (
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Ray starting from an origin point and extending to infinity in a given direction.
type Ray struct {
	// Origin is the starting point of the ray
	Origin vec3.T

	// Dir is ray direction as a unit vector
	Dir vec3.T
}
