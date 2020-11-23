package physics

import "github.com/johanhenriksson/goworld/math/vec3"

// Object is the generic interface for physics objects
type Object interface {
	Intersect(*Ray) (bool, vec3.T)
}
