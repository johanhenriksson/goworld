package collider

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/physics"
	"github.com/johanhenriksson/goworld/math/vec3"
)

type Box struct {
	Center vec3.T
	Size   vec3.T
}

type box struct {
	object.T
	args  Box
	shape physics.Box
}

func NewBox(args Box) T {
	half := args.Size.Scaled(0.5)
	return object.New(&box{
		args: args,
		shape: physics.Box{
			Min: args.Center.Sub(half),
			Max: args.Center.Add(half),
		},
	})
}

func (s *box) Intersect(ray *physics.Ray) (bool, vec3.T) {
	return s.shape.Intersect(ray)
}

func (s *box) Update(scene object.T, dt float32) {
	// this isnt enough for rotation
	sz := s.args.Size.Scaled(0.5)
	s.shape.Min = s.Transform().Project(s.args.Center.Sub(sz))
	s.shape.Max = s.Transform().Project(s.args.Center.Add(sz))
}
