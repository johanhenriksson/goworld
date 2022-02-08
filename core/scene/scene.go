package scene

import (
	"github.com/johanhenriksson/goworld/core/object"
)

type T interface {
	object.T
}

// Scene graph root
type scene struct {
	object.T
}

// NewScene creates a new scene.
func New() T {
	return &scene{
		T: object.New("Scene"),
	}
}
