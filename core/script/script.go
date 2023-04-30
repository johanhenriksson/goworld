package script

import "github.com/johanhenriksson/goworld/core/object"

type T interface {
	object.T
}

type script struct {
	object.T
	fn Behavior
}

type Behavior func(scene, self object.T, dt float32)

func New(fn Behavior) T {
	return object.New(&script{
		fn: fn,
	})
}

func (s *script) Update(scene object.T, dt float32) {
	s.fn(scene, s, dt)
}
