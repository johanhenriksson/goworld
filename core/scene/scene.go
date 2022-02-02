package scene

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
)

type T interface {
	object.T

	Camera() camera.T
	SetCamera(camera.T)
}

// Scene graph root
type scene struct {
	object.T

	// Active camera
	camera camera.T
}

// NewScene creates a new scene.
func New() T {
	return &scene{
		T:      object.New("Scene"),
		camera: nil,
	}
}

func (s *scene) Camera() camera.T {
	return s.camera
}

func (s *scene) SetCamera(cam camera.T) {
	s.camera = cam
}
