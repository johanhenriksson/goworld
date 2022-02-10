package physics

import (
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Plane is a very flat, infinite shape
type Plane struct {
	Normal vec3.T
	D      float32
}

// Intersect checks if a ray intersects the plane.
func (p *Plane) Intersect(ray *Ray) (bool, float32, vec3.T) {
	denom := vec3.Dot(ray.Dir, p.Normal)
	if denom > 0 {
		t := -(vec3.Dot(ray.Origin, p.Normal) + p.D) / denom
		if t > 0 {
			return true, t, ray.Origin.Add(ray.Dir.Scaled(t))
		}
	}
	return false, 0.0, vec3.Zero
}
