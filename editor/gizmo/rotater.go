package gizmo

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
)

// Rotater Gizmo can be used to reposition objects in the 3D scene.
type Rotater struct {
	object.Object

	target transform.T

	Sphere    *lines.Sphere
	Rigidbody *physics.RigidBody
	Shape     *physics.Sphere

	// screen size scaling factor
	size float32

	initialRot quat.T
	axisFilter vec3.T
	start      vec3.T

	eye      vec3.T
	viewport draw.Viewport
	vp       mat4.T
	fov      float32
	scale    float32
	dragging bool
}

var _ Gizmo = &Rotater{}

// NewRotater creates a new mover gizmo
func NewRotater(pool object.Pool) *Rotater {
	g := object.New(pool, "Rotater Gizmo", &Rotater{
		size: 0.1,

		Sphere: lines.NewSphere(pool, lines.SphereArgs{
			Radius: 1,
			Color:  color.Purple,
		}),
		Shape:     physics.NewSphere(pool, 1),
		Rigidbody: physics.NewRigidBody(pool, 0),
	})

	g.Rigidbody.Layer.Set(2)
	g.Sphere.SetAxisColors(color.Red, color.Green, color.Blue)

	return g
}

func (g *Rotater) Target() transform.T {
	return g.target
}

func (g *Rotater) SetTarget(t transform.T) {
	g.Transform().SetParent(t)
	g.target = t
	if t == nil {
		g.DragEnd(mouse.NopEvent(), physics.RaycastHit{})
	}
}

func (g *Rotater) CanDeselect() bool {
	return true
}

func (g *Rotater) DragStart(e mouse.Event, hit physics.RaycastHit) {
	dir := g.Transform().Unproject(hit.Point).Normalized()

	filter := vec3.One
	if math.Abs(dir.X) < 0.1 {
		filter.X = 0
	}
	if math.Abs(dir.Y) < 0.1 {
		filter.Y = 0
	}
	if math.Abs(dir.Z) < 0.1 {
		filter.Z = 0
	}
	g.axisFilter = filter

	g.dragging = true
	g.start = dir.Mul(filter).Normalized()
	g.initialRot = g.target.Rotation()
}

func (g *Rotater) DragEnd(e mouse.Event, hit physics.RaycastHit) {
	g.dragging = false
}

func (g *Rotater) DragMove(e mouse.Event, hit physics.RaycastHit) {
	if hit.Shape == nil {
		return
	}
	dir := g.Transform().Unproject(hit.Point).Normalized()
	dir = dir.Mul(g.axisFilter).Normalized()

	rot := quat.BetweenVectors(g.start, dir)

	g.target.SetRotation(rot.Mul(g.initialRot))
}

func (g *Rotater) Hover(hovering bool, shape physics.Shape) {

}

func (g *Rotater) PreDraw(args draw.Args, scene object.Object) error {
	g.eye = args.Camera.Position
	g.fov = args.Camera.Fov
	g.vp = args.Camera.ViewProj
	g.viewport = args.Camera.Viewport
	return nil
}

func (g *Rotater) Update(scene object.Component, dt float32) {
	g.Object.Update(scene, dt)

	// the gizmo should be displayed at the same size irrespectively of its distance to the camera.
	// we can undo the effects of perspective projection by measuring how much a vector would be "squeezed"
	// at the current distance form the camera, and then applying a scaling factor to counteract it.
	distance := vec3.Distance(g.eye, g.Transform().WorldPosition())
	worldSize := (2 * math.Tan(math.DegToRad(g.fov)/2.0)) * distance
	f := g.size * worldSize

	g.scale = f
	g.Transform().SetWorldScale(vec3.New(f, f, f))
	g.Transform().SetWorldRotation(quat.Ident())
}

func (g *Rotater) Dragging() bool          { return g.dragging }
func (g *Rotater) Viewport() draw.Viewport { return g.viewport }
func (g *Rotater) Camera() mat4.T          { return g.vp }

func (m *Rotater) ToolMouseEvent(e mouse.Event, hover physics.RaycastHit) {
	HandleMouse(m, e, hover)
}
