package engine

import (
	"github.com/johanhenriksson/goworld/math/mat4"
)

// Scene graph root
type Scene struct {
	// Active camera
	Camera *Camera

	// Root Objects
	Objects []Component

	// List of all lights in the scene
	Lights []Light
}

// NewScene creates a new scene.
func NewScene() *Scene {
	return &Scene{
		Camera:  nil,
		Objects: []Component{},
		Lights:  []Light{},
	}
}

// Add an object to the scene
func (s *Scene) Add(object Component) {
	// TODO: keep track of lights
	s.Objects = append(s.Objects, object)
}

// Update the scene.
func (s *Scene) Update(dt float32) {
	if s.Camera != nil {
		// update camera first
		s.Camera.Update(dt)
	}

	// update root objects
	for _, obj := range s.Objects {
		obj.Update(dt)
	}
}

func (s *Scene) CollectWithArgs(pass DrawPass, args DrawArgs) {
	for _, obj := range s.Objects {
		obj.Collect(pass, args)
	}
}

func (s *Scene) Collect(pass DrawPass) {
	p := s.Camera.Projection
	v := s.Camera.View
	m := mat4.Ident()
	vp := p.Mul(&v)

	/* DrawArgs will be copied down recursively into the scene graph.
	 * Each object adds its transformation matrix before passing
	 * it on to their children */
	args := DrawArgs{
		Projection: p,
		View:       v,
		VP:         vp,
		MVP:        vp,
		Transform:  m,
		Position:   s.Camera.Position,
		Pass:       pass.Type(),
	}

	s.CollectWithArgs(pass, args)
}
