package engine

import (
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render"

	mgl "github.com/go-gl/mathgl/mgl32"
)

/* Scene Graph */
type Scene struct {
	/* Active camera */
	Camera *Camera

	/* Root Objects */
	Objects []*Object

	World *physics.World

	/* temporary: list of all lights in the scene */
	Lights []Light
}

func NewScene() *Scene {
	s := &Scene{
		Camera:  nil,
		Objects: []*Object{},
		World:   physics.NewWorld(),
		Lights: []Light{},
	}
	
	return s
}

func (s *Scene) Add(object *Object) {
	/* TODO look for lights - maybe not here? */
	s.Objects = append(s.Objects, object)
}

func (s *Scene) Draw(pass string, shader *render.ShaderProgram) {
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

		Pass:   pass,
		Shader: shader,
	}

	s.DrawCall(args)
}

func (s *Scene) DrawCall(args render.DrawArgs) {
	/* draw root objects */
	for _, obj := range s.Objects {
		obj.Draw(args)
	}
}

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

	/* physics step */
	s.World.Update()
}