package shape

import "github.com/johanhenriksson/goworld/math/vec3"

type Sphere struct {
	Center vec3.T
	Radius float32
}

func (s *Sphere) IntersectsSphere(other *Sphere) bool {
	sepAxis := s.Center.Sub(other.Center)
	radiiSum := s.Radius + other.Radius
	intersects := sepAxis.LengthSqr() < (radiiSum * radiiSum)
	return intersects
}
