package physics

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Sphere is a 3D circle. Yeah.
type Sphere struct {
	// Center is the midpoint of the sphere
	Center vec3.T

	// Radius of the sphere
	Radius float32
}

// Intersect checks a ray for intersection with the sphere.
func (s *Sphere) Intersect(ray *Ray) (bool, vec3.T) {
	oc := ray.Origin.Sub(s.Center)
	b := vec3.Dot(ray.Dir, oc)
	c := oc.LengthSqr() - s.Radius*s.Radius

	b2c := b*b - c
	if b2c >= 0 {
		// hit!
		t1 := -b - math.Sqrt(b2c)
		return true, ray.Origin.Add(ray.Dir.Scaled(t1))
	}

	return false, vec3.Zero
}
