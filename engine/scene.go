package engine

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
)

// Scene graph root
type Scene struct {
	object.T

	// Active camera
	Camera camera.T

	// List of all lights in the scene
	Lights []Light
}

// NewScene creates a new scene.
func NewScene() *Scene {
	return &Scene{
		T:      object.New("Scene"),
		Camera: nil,
		Lights: []Light{},
	}
}

// Update the scene.
func (s *Scene) Update(dt float32) {
	if s.Camera != nil {
		// update camera first
		s.Camera.Update(dt)
	}
	s.T.Update(dt)
}
