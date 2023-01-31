package collider

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/physics"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type Selectable interface {
	object.T
	Select(mouse.Event, T)
	Deselect(mouse.Event)
	SelectedMouseEvent(mouse.Event)
}

type Manager interface {
	object.T
	mouse.Handler
}

type manager struct {
	object.T
	scene    object.T
	selected Selectable

	camera   mat4.T
	viewport vec2.T
	eye      vec3.T
}

func NewManager() Manager {
	return object.New(&manager{})
}

func (m *manager) MouseEvent(e mouse.Event) {
	if m.scene == nil {
		return
	}

	vpi := m.camera.Invert()
	cursor := e.Position().Div(m.viewport).Sub(vec2.New(0.5, 0.5)).Scaled(2)

	if m.selected != nil {
		if e.Action() == mouse.Release {
			m.selected.Deselect(e)
			m.selected = nil
		} else {
			m.selected.SelectedMouseEvent(e)
		}
	} else if e.Action() == mouse.Press {
		near := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 1))
		far := vpi.TransformPoint(vec3.New(cursor.X, cursor.Y, 0))
		dir := far.Sub(near).Normalized()

		// return closest hit
		colliders := object.Query[T]().Collect(m.scene)
		var closest T
		closestDist := float32(math.InfPos)
		for _, collider := range colliders {
			hit, point := collider.Intersect(&physics.Ray{
				Origin: near,
				Dir:    dir,
			})
			if hit {
				dist := vec3.Distance(point, m.eye)
				if dist < closestDist {
					closest = collider
					closestDist = dist
				}
			}
		}
		if closest != nil {
			if selectable, ok := object.GetInParents[Selectable](closest); ok {
				m.selected = selectable
				selectable.Select(e, closest)
			}
		} else if m.selected != nil {
			// we hit nothing, deselect
			m.selected.Deselect(e)
		}
	}
}

func (m *manager) PreDraw(args render.Args, scene object.T) error {
	m.scene = scene
	m.camera = args.VP
	m.viewport = vec2.New(float32(args.Viewport.Width), float32(args.Viewport.Height))
	m.eye = args.Position
	return nil
}
