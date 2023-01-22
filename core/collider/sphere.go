package collider

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/physics"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type T interface {
	object.Component

	Intersect(ray *physics.Ray) (bool, vec3.T)
}

type Sphere struct {
	Center vec3.T
	Radius float32
}

type sphere struct {
	object.Component
	args  Sphere
	shape physics.Sphere
}

func NewSphere(args Sphere) T {
	return &sphere{
		Component: object.NewComponent(),
		args:      args,
		shape: physics.Sphere{
			Center: args.Center,
			Radius: args.Radius,
		},
	}
}

func (s *sphere) Intersect(ray *physics.Ray) (bool, vec3.T) {
	return s.shape.Intersect(ray)
}

func (s *sphere) Update(dt float32) {
	s.shape.Center = s.Transform().Project(s.args.Center)
}
