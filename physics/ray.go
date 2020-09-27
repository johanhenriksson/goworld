package physics

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Ray struct {
	Origin vec3.T
	Dir    vec3.T
}

func (ray Ray) IntersectBox(box *Box) (bool, vec3.T) {
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

	return true, vec3.T{hit[0], hit[1], hit[2]}
}

func (ray Ray) IntersectPlane(p *Plane) (bool, float32, vec3.T) {
	denom := vec3.Dot(ray.Dir, p.Normal)
	if denom > 0 {
		t := -(vec3.Dot(ray.Origin, p.Normal) + p.D) / denom
		if t > 0 {
			return true, t, ray.Origin.Add(ray.Dir.Scaled(t))
		}
	}
	return false, 0.0, vec3.Zero
}

func (ray Ray) IntersectSphere(s *Sphere) (bool, float32, vec3.T) {
	oc := ray.Origin.Sub(s.Center)
	b := vec3.Dot(ray.Dir, oc)
	c := oc.LengthSqr() - s.Radius*s.Radius

	b2c := b*b - c
	if b2c >= 0 {
		// hit!
		t1 := -b - math.Sqrt(b2c)
		return true, t1, ray.Origin.Add(ray.Dir.Scaled(t1))
	}

	return false, 0.0, vec3.Zero
}
