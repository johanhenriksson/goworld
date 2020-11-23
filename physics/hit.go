package physics

import "github.com/johanhenriksson/goworld/math/vec3"

// Hit represents a raycast hit
type Hit struct {
	Point  vec3.T
	Object Object
}
