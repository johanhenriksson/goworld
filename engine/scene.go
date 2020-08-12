package engine

import (
	"github.com/johanhenriksson/goworld/render"

	mgl "github.com/go-gl/mathgl/mgl32"
)

// Scene graph root
type Scene struct {
	/* Active camera */
	Camera *Camera

	/* Root Objects */
	Objects []*Object

	/* temporary: list of all lights in the scene */
	Lights []Light
}

// NewScene creates a new scene.
func NewScene() *Scene {
	s := &Scene{
		Camera:  nil,
		Objects: []*Object{},
		Lights:  []Light{},
	}

	return s
}

// Add an object to the scene
func (s *Scene) Add(object *Object) {
	/* TODO look for lights - maybe not here? */
	s.Objects = append(s.Objects, object)
}

// DrawPass draws the scene using the default camera and a specific render pass
func (s *Scene) DrawPass(pass render.DrawPass) {
	if s.Camera == nil {
		return
	}

	p := s.Camera.Projection
	v := s.Camera.View
	m := mgl.Ident4()
	vp := p.Mul4(v)
	// mvp := vp * m

	/* DrawArgs will be copied down recursively into the scene graph.
	 * Each object adds its transformation matrix before passing
	 * it on to their children */
	args := render.DrawArgs{
		Projection: p,
		View:       v,
		VP:         vp,
		MVP:        vp,
		Transform:  m,

		Pass: pass,
	}

	s.Draw(args)
}

// Draw the scene using the provided render arguments
func (s *Scene) Draw(args render.DrawArgs) {
	/* draw root objects */
	for _, obj := range s.Objects {
		obj.Draw(args)
	}
}

// Update the scene.
func (s *Scene) Update(dt float32) {
	if s.Camera != nil {
		/* update camera first */
		s.Camera.Update(dt)
	}

	/* update root objects */
	for _, obj := range s.Objects {
		obj.Update(dt)
	}

	/* test: position first light on camera */
	//s.lights[0].Position = s.Camera.Position
}
