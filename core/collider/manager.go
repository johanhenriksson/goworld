package collider

import (
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
	mouse.Handler
}

type manager struct {
	object.Component
	scene    object.T
	camera   mat4.T
	viewport vec2.T
	forward  vec3.T

	drag       object.T
	axis       vec3.T
	screenAxis vec2.T
	start      vec2.T
}

func NewManager() Manager {
	return &manager{
		Component: object.NewComponent(),
	}
}

func (m *manager) MouseEvent(e mouse.Event) {
	if m.scene == nil {
		return
	}

	vpi := m.camera.Invert()
	cursor := e.Position().Div(m.viewport).Sub(vec2.New(0.5, 0.5)).Scaled(2)

	if e.Action() == mouse.Release {
		m.drag = nil
		return
	} else if e.Action() == mouse.Move {
		if m.drag != nil {
			delta := m.start.Sub(cursor)
			mag := -5 * vec2.Dot(delta, m.screenAxis) / m.screenAxis.Length()
			m.start = cursor
			pos := m.drag.Transform().Position().Add(m.axis.Scaled(mag))
			m.drag.Transform().SetPosition(pos)
			e.Consume()
		}
		return
	}

	near := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 1))
	far := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 0))
	dir := far.Sub(near).Normalized()

	// todo: return closest hit
	colliders := query.New[T]().Collect(m.scene)
	for _, collider := range colliders {
		hit, _ := collider.Intersect(&physics.Ray{
			Origin: near,
			Dir:    dir,
		})
		if hit {
			if e.Action() == mouse.Press {
				m.start = cursor
				m.drag = collider.Object().Parent()
				axisName := collider.Object().Name()[:1]
				switch axisName {
				case "X":
					m.axis = vec3.UnitX
				case "Y":
					m.axis = vec3.UnitY
				case "Z":
					m.axis = vec3.UnitZ
				default:
					return
				}

				localDir := m.drag.Transform().ProjectDir(m.axis)
				m.screenAxis = m.camera.TransformDir(localDir).XY().Normalized()

				e.Consume()
				break
			}
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
