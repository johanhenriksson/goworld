package script

import "github.com/johanhenriksson/goworld/core/object"

type T interface {
	object.Component
}

type script struct {
	object.Component
	fn Behavior
}

type Behavior func(scene, self object.Component, dt float32)

func New(pool object.Pool, fn Behavior) T {
	return object.NewComponent(pool, &script{
		fn: fn,
	})
}

func (s *script) Update(scene object.Component, dt float32) {
	s.fn(scene, s, dt)
}
