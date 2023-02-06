package shape

import (
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Plane struct {
	Normal   vec3.T
	Distance float32
}

func (p *Plane) normalize() {
	length := p.Normal.LengthSqr() + p.Distance*p.Distance
	p.Normal = p.Normal.Scaled(1 / length)
	p.Distance /= length
}

func (p *Plane) DistanceToPoint(point vec3.T) float32 {
	return vec3.Dot(p.Normal, point) + p.Distance
}

type Frustum struct {
	Front, Back, Left, Right, Top, Bottom Plane
}

func (f *Frustum) IntersectsSphere(s *Sphere) bool {
	if f.Left.DistanceToPoint(s.Center) <= -s.Radius {
		return false
	}
	if f.Right.DistanceToPoint(s.Center) <= -s.Radius {
		return false
	}
	if f.Top.DistanceToPoint(s.Center) <= -s.Radius {
		return false
	}
	if f.Bottom.DistanceToPoint(s.Center) <= -s.Radius {
		return false
	}
	if f.Front.DistanceToPoint(s.Center) <= -s.Radius {
		return false
	}
	if f.Back.DistanceToPoint(s.Center) <= -s.Radius {
		return false
	}
	return true
}

func FrustumFromMatrix(vp mat4.T) Frustum {
	f := Frustum{
		Left: Plane{
			Normal: vec3.T{
				X: vp[0+3] + vp[0+0],
				Y: vp[4+3] + vp[4+0],
				Z: vp[8+3] + vp[8+0],
			},
			Distance: vp[12+3] + vp[12+0],
		},
		Right: Plane{
			Normal: vec3.T{
				X: vp[0+3] - vp[0+0],
				Y: vp[4+3] - vp[4+0],
				Z: vp[8+3] - vp[8+0],
			},
			Distance: vp[12+3] - vp[12+0],
		},
		Top: Plane{
			Normal: vec3.T{
				X: vp[0+3] - vp[0+1],
				Y: vp[4+3] - vp[4+1],
				Z: vp[8+3] - vp[8+1],
			},
			Distance: vp[12+3] - vp[12+1],
		},
		Bottom: Plane{
			Normal: vec3.T{
				X: vp[0+3] + vp[0+1],
				Y: vp[4+3] + vp[4+1],
				Z: vp[8+3] + vp[8+1],
			},
			Distance: vp[12+3] + vp[12+1],
		},
		Back: Plane{
			Normal: vec3.T{
				X: vp[0+3] + vp[0+2],
				Y: vp[4+3] + vp[4+2],
				Z: vp[8+3] + vp[8+2],
			},
			Distance: vp[12+3] + vp[12+2],
		},
		Front: Plane{
			Normal: vec3.T{
				X: vp[0+3] - vp[0+2],
				Y: vp[4+3] - vp[4+2],
				Z: vp[8+3] - vp[8+2],
			},
			Distance: vp[12+3] - vp[12+2],
		},
	}
	f.Front.normalize()
	f.Back.normalize()
	f.Top.normalize()
	f.Bottom.normalize()
	f.Left.normalize()
	f.Right.normalize()
	return f
}
