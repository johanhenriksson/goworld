package gizmo

import (
	"github.com/johanhenriksson/goworld/core/draw"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	. "github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
)

// Mover Gizmo can be used to reposition objects in the 3D scene.
type Mover struct {
	Object

	target transform.T

	X *Arrow
	Y *Arrow
	Z *Arrow

	XY *plane.Plane
	XZ *plane.Plane
	YZ *plane.Plane

	// screen size scaling factor
	size float32

	eye         vec3.T
	axis        vec3.T
	screenAxis  vec2.T
	start       vec2.T
	viewport    draw.Viewport
	vp          mat4.T
	fov         float32
	scale       float32
	sensitivity float32
	dragging    bool
	hoverScale  vec3.T
}

var _ Gizmo = &Mover{}

// NewMover creates a new mover gizmo
func NewMover(pool Pool) *Mover {
	side := float32(0.4)
	planeAlpha := float32(0.33)

	s := side * 0.5

	g := NewObject(pool, "Mover Gizmo", &Mover{
		size:        0.125,
		sensitivity: 6,
		hoverScale:  vec3.New(1.1, 1.1, 1.1),

		X: Builder(NewArrow(pool, color.Red)).
			Rotation(quat.Euler(0, 0, 270)).
			Create(),

		Y: NewArrow(pool, color.Green),

		Z: Builder(NewArrow(pool, color.Blue)).
			Rotation(quat.Euler(90, 0, 0)).
			Create(),

		// XY Plane
		XY: Builder(plane.New(pool, plane.Args{
			Size: vec2.New(side, side),
			Mat:  material.TransparentForward(),
		})).
			Position(vec3.New(s, s, 0)).
			Rotation(quat.Euler(90, 0, 0)).
			Texture(texture.Diffuse, color.Blue.WithAlpha(planeAlpha)).
			Create(),

		// XZ Plane
		XZ: Builder(plane.New(pool, plane.Args{
			Size: vec2.New(side, side),
			Mat:  material.TransparentForward(),
		})).
			Rotation(quat.Euler(0, 90, 0)).
			Position(vec3.New(s, 0, s)).
			Texture(texture.Diffuse, color.Green.WithAlpha(planeAlpha)).
			Create(),

		// YZ Plane
		YZ: Builder(plane.New(pool, plane.Args{
			Size: vec2.New(side, side),
			Mat:  material.TransparentForward(),
		})).
			Position(vec3.New(0, s, s)).
			Rotation(quat.Euler(0, 0, 90)).
			Texture(texture.Diffuse, color.Red.WithAlpha(planeAlpha)).
			Create(),
	})

	return g
}

func (g *Mover) Target() transform.T {
	return g.target
}

func (g *Mover) SetTarget(t transform.T) {
	g.Transform().SetParent(t)
	g.target = t
	if t == nil {
		g.DragEnd(mouse.NopEvent(), physics.RaycastHit{})
	}
}

func (g *Mover) CanDeselect() bool {
	return true
}

func (g *Mover) getColliderAxis(collider physics.Shape) vec3.T {
	axisObj := collider.Parent()
	switch axisObj {
	case g.X:
		return vec3.UnitX
	case g.Y:
		return vec3.UnitY
	case g.Z:
		return vec3.UnitZ
	}
	return vec3.Zero
}

func (g *Mover) DragStart(e mouse.Event, hit physics.RaycastHit) {
	g.dragging = true

	g.axis = g.getColliderAxis(hit.Shape)
	cursor := g.viewport.NormalizeCursor(e.Position())
	g.start = cursor

	localDir := g.Transform().ProjectDir(g.axis)
	g.screenAxis = g.vp.TransformDir(localDir).XY().Normalized()
}

func (g *Mover) DragEnd(e mouse.Event, hit physics.RaycastHit) {
	g.dragging = false
}

func (g *Mover) DragMove(e mouse.Event, hit physics.RaycastHit) {
	if e.Action() == mouse.Move {
		cursor := g.viewport.NormalizeCursor(e.Position())
		axisLen := g.screenAxis.Length()
		if axisLen == 0 {
			return
		}

		delta := g.start.Sub(cursor)
		mag := -1 * g.sensitivity * g.scale * vec2.Dot(delta, g.screenAxis) / axisLen
		g.start = cursor
		pos := g.Transform().WorldPosition().Add(g.axis.Scaled(mag))

		if g.target != nil {
			g.target.SetWorldPosition(pos)
		}
	}
}

func (g *Mover) Hover(hovering bool, shape physics.Shape) {
	if hovering {
		// hover start
		axis := g.getColliderAxis(shape)
		switch axis {
		case vec3.UnitX:
			g.X.Hover.Set(true)
			g.Y.Hover.Set(false)
			g.Z.Hover.Set(false)
		case vec3.UnitY:
			g.X.Hover.Set(false)
			g.Y.Hover.Set(true)
			g.Z.Hover.Set(false)
		case vec3.UnitZ:
			g.X.Hover.Set(false)
			g.Y.Hover.Set(false)
			g.Z.Hover.Set(true)
		}
	}
	if !hovering {
		g.X.Hover.Set(false)
		g.Y.Hover.Set(false)
		g.Z.Hover.Set(false)
	}
}

func (g *Mover) PreDraw(args draw.Args, scene Object) error {
	g.eye = args.Camera.Position
	g.fov = args.Camera.Fov
	g.vp = args.Camera.ViewProj
	g.viewport = args.Camera.Viewport
	return nil
}

func (g *Mover) Update(scene Component, dt float32) {
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

func (g *Mover) Dragging() bool          { return g.dragging }
func (g *Mover) Viewport() draw.Viewport { return g.viewport }
func (g *Mover) Camera() mat4.T          { return g.vp }

func (m *Mover) ToolMouseEvent(e mouse.Event, hover physics.RaycastHit) {
	HandleMouse(m, e, hover)
}
