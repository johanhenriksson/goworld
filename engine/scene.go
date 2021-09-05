package engine

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
)

// Scene graph root
type Scene struct {
	object.T

	// Active camera
	Camera camera.T

	// List of all lights in the scene
	Lights []render.Light
}

// NewScene creates a new scene.
func NewScene() *Scene {
	return &Scene{
		T:      object.New("Scene"),
		Camera: nil,
		Lights: []render.Light{},
	}
}
