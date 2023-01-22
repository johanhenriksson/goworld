package physics

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Box holds an axis-aligned bounding box.
type Box struct {
	Min vec3.T
	Max vec3.T
}

func NewBox(center, extents vec3.T) Box {
	half := extents.Scaled(0.5)
	return Box{
		Min: center.Sub(half),
		Max: center.Add(half),
	}
}

func (box Box) Center() vec3.T {
	return box.Min.Add(box.Max).Scaled(0.5)
}

func (box Box) Extents() vec3.T {
	return box.Max.Sub(box.Min)
}

func (box Box) Corners() [8]vec3.T {
	return [8]vec3.T{
		box.Min, // 000
		vec3.New(box.Min.X, box.Min.Y, box.Max.Z), // 001
		vec3.New(box.Min.X, box.Max.Y, box.Min.Z), // 010
		vec3.New(box.Min.X, box.Max.Y, box.Max.Z), // 011
		vec3.New(box.Max.X, box.Min.Y, box.Min.Z), // 100
		vec3.New(box.Max.X, box.Min.Y, box.Max.Z), // 101
		vec3.New(box.Max.X, box.Max.Y, box.Min.Z), // 110
		box.Max, // 111
	}
}

// Intersect a ray with this box. Returns a value indicating if it hit, and if so, the point of intersection.
func (box *Box) Intersect(ray *Ray) (bool, vec3.T) {
	mins := box.Min.Sub(ray.Origin).Div(ray.Dir)
	maxs := box.Max.Sub(ray.Origin).Div(ray.Dir)
	tMin := math.Max(math.Max(math.Min(mins.X, maxs.X), math.Min(mins.Y, maxs.Y)), math.Min(mins.Z, maxs.Z))
	tMax := math.Min(math.Min(math.Max(mins.X, maxs.X), math.Max(mins.Y, maxs.Y)), math.Max(mins.Z, maxs.Z))

	// if tmax < 0, ray (line) is intersecting AABB, but whole AABB is behing us
	if tMax < 0 {
		return false, vec3.Zero
	}

	// if tmin > tmax, ray doesn't intersect AABB
	if tMin > tMax {
		return false, vec3.Zero
	}

	if tMin < 0 {
		return true, ray.Origin.Add(ray.Dir.Scaled(tMax))
	}
	return true, ray.Origin.Add(ray.Dir.Scaled(tMin))
}
