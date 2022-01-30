package scene

import (
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/light"
)

type T interface {
	object.T

	Camera() camera.T
	SetCamera(camera.T)

	Lights() []light.T
}

// Scene graph root
type scene struct {
	object.T

	// Active camera
	camera camera.T

	// List of all lights in the scene
	lights []light.T
}

// NewScene creates a new scene.
func New() T {
	return &scene{
		T:      object.New("Scene"),
		camera: nil,
		lights: []light.T{
			{ // directional light
				Intensity:  0.8,
				Color:      color.RGB(0.9*0.973, 0.9*0.945, 0.9*0.776),
				Type:       light.Directional,
				Projection: mat4.Orthographic(-71, 120, -20, 140, -10, 140),
				Position:   vec3.New(-2, 2, -1),
				Shadows:    false,
			},
			{ // point light
				Intensity:   1.8,
				Range:       10.0,
				Color:       color.Yellow,
				Type:        light.Point,
				Position:    vec3.New(12, 7, 12),
				Shadows:     false,
				Attenuation: light.DefaultAttenuation,
			},
		},
	}
}

func (s *scene) Lights() []light.T {
	return s.lights
}

func (s *scene) Camera() camera.T {
	return s.camera
}

func (s *scene) SetCamera(cam camera.T) {
	s.camera = cam
}
