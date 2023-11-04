package server

import "github.com/johanhenriksson/goworld/math/vec3"

type Identity uint64

type Entity interface {
	ID() Identity
	Name() string
	Position() vec3.T
	SetPosition(vec3.T)
}
