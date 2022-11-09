package collider

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/physics"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type Manager interface {
	object.Component
}

type manager struct {
	object.Component
	scene    object.T
	camera   mat4.T
	viewport vec2.T
	forward  vec3.T
}

func NewManager() Manager {
	return &manager{
		Component: object.NewComponent(),
	}
}

func (m *manager) MouseEvent(e mouse.Event) {
	if e.Action() != mouse.Press {
		return
	}
	if m.scene == nil {
		return
	}

	vpi := m.camera.Invert()
	cursor := e.Position().Div(m.viewport).Sub(vec2.New(0.5, 0.5)).Scaled(2)

	near := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 1))
	far := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 0))
	dir := far.Sub(near).Normalized()

	colliders := query.New[T]().Collect(m.scene)
	for _, collider := range colliders {
		hit, point := collider.Intersect(&physics.Ray{
			Origin: near,
			Dir:    dir,
		})
		if hit {
			fmt.Println("hit", collider.Object().Name(), "at", point)
			e.Consume()
		}
	}
}

func (m *manager) PreDraw(args render.Args, scene object.T) error {
	m.scene = scene
	m.camera = args.VP
	m.viewport = vec2.New(float32(args.Viewport.Width), float32(args.Viewport.Height))
	m.forward = args.Forward
	return nil
}
