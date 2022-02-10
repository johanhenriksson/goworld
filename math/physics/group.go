package physics

import (
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/vec3"
)

// Group physics objects together
type Group struct {
	objects []Object
}

// NewGroup creates a new physics group from a set of physics objects.
func NewGroup(objects ...Object) *Group {
	return &Group{objects}
}

// Add a new object to the group.
func (g *Group) Add(object Object) {
	g.objects = append(g.objects)
}

// Intersect performs ray intersection with every object in the group
// and returns a list of all intersecting objects.
func (g *Group) Intersect(ray *Ray) []Hit {
	hits := make([]Hit, 0, 8)
	for _, object := range g.objects {
		hit, point := object.Intersect(ray)
		if hit {
			hits = append(hits, Hit{
				Point:  point,
				Object: object,
			})
		}
	}
	return hits
}

// ClosestIntersect performs a ray intersection with the group
// and returns the closest intersecting object.
func (g *Group) ClosestIntersect(ray *Ray) (bool, Hit) {
	hits := g.Intersect(ray)
	best := Hit{}
	bestDist := math.MaxValue
	if len(hits) == 0 {
		return false, best
	}
	for _, hit := range hits {
		dist := vec3.Distance(hit.Point, ray.Origin)
		if dist < bestDist {
			bestDist = dist
			best = hit
		}
	}
	return true, best
}
