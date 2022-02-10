package physics

import (
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Box holds a centered an axis-aligned bounding box.
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
	// Fast Ray-Box Intersection by Andrew Woo
	// from "Graphics Gems", Academic Press, 1990

	const (
		RIGHT  = 0
		LEFT   = 1
		MIDDLE = 2
		DIM    = 3
	)

	hit := [DIM]float32{}
	minB := box.Min.Slice()
	maxB := box.Max.Slice()
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
