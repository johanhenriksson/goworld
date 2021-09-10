package scene

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type T interface {
	object.T

	Camera() camera.T
	SetCamera(camera.T)

	Lights() []render.Light
}

// Scene graph root
type scene struct {
	object.T

	// Active camera
	camera camera.T

	// List of all lights in the scene
	lights []render.Light
}

// NewScene creates a new scene.
func New() T {
	return &scene{
		T:      object.New("Scene"),
		camera: nil,
		lights: []render.Light{
			{ // directional light
				Intensity:  1.6,
				Color:      vec3.New(0.9*0.973, 0.9*0.945, 0.9*0.776),
				Type:       render.DirectionalLight,
				Projection: mat4.Orthographic(-71, 120, -20, 140, -10, 140),
				Position:   vec3.New(-2, 2, -1),
				Shadows:    false,
			},
		},
	}
}

func (s *scene) Lights() []render.Light {
	return s.lights
}

func (s *scene) Camera() camera.T {
	return s.camera
}

func (s *scene) SetCamera(cam camera.T) {
	s.camera = cam
}