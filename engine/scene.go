package engine

import "github.com/johanhenriksson/goworld/engine/object"

// Scene graph root
type Scene struct {
	object.T

	// Active camera
	Camera *Camera

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
