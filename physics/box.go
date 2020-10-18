package physics

import (
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Box holds a centered an axis-aligned bounding box.
type Box struct {
	Extents vec3.T
	Center  vec3.T
}

// Min returns the lower point
func (box Box) Min() vec3.T {
	return vec3.T{
		X: box.Center.X - box.Extents.X/2,
		Y: box.Center.Y - box.Extents.Y/2,
		Z: box.Center.Z - box.Extents.Z/2,
	}
}

// Max returns the upper point
func (box Box) Max() vec3.T {
	return vec3.T{
		X: box.Center.X + box.Extents.X/2,
		Y: box.Center.Y + box.Extents.Y/2,
		Z: box.Center.Z + box.Extents.Z/2,
	}
}

// Intersect a ray with this box. Returns a value indicating if it hit, and if so, the point of intersection.
func (box *Box) Intersect(ray *Ray) (bool, vec3.T) {
	// Fast Ray-Box Intersection by Andrew Woo
	// from "Graphics Gems", Academic Press, 1990

	const (
		RIGHT  = 0
		LEFT   = 1
		MIDDLE = 2
		DIM    = 3
	)

	hit := [DIM]float32{}
	minB := box.Min().Slice()
	maxB := box.Max().Slice()
	inside := true
	maxT := [DIM]float32{}
	origin := ray.Origin.Slice()
	dir := ray.Dir.Slice()
	candidate := [DIM]float32{}
	quadrant := [DIM]uint8{}
	whichPlane := 0

	// Find candidate planes
	for i := 0; i < DIM; i++ {
		if origin[i] < minB[i] {
			quadrant[i] = LEFT
			candidate[i] = minB[i]
			inside = false
		} else if origin[i] > maxB[i] {
			quadrant[i] = RIGHT
			candidate[i] = maxB[i]
			inside = false
		} else {
			quadrant[i] = MIDDLE
		}
	}

	// ray origin is inside the bounding box
	if inside {
		return true, ray.Origin
	}

	// calculate T distance to candidate planes
	for i := 0; i < DIM; i++ {
		if quadrant[i] != MIDDLE && dir[i] != 0 {
			maxT[i] = (candidate[i] - origin[i]) / dir[i]
		} else {
			maxT[i] = -1
		}
	}

	// choose largest maxT
	for i := 0; i < DIM; i++ {
		if maxT[whichPlane] < maxT[i] {
			whichPlane = i
		}
	}

	// make sure final candidate is actually inside the bounding box
	if maxT[whichPlane] < 0 {
		return false, vec3.Zero
	}
	for i := 0; i < DIM; i++ {
		if whichPlane != i {
			hit[i] = origin[i] - maxT[whichPlane]*dir[i]
			if hit[i] < minB[i] || hit[i] > maxB[i] {
				return false, vec3.Zero
			}
		} else {
			hit[i] = candidate[i]
		}
	}

	return true, vec3.New(hit[0], hit[1], hit[2])
}
