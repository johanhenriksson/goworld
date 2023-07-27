package gizmo

import (
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/physics"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Mover Gizmo can be used to reposition objects in the 3D scene.
type Mover struct {
	object.Object

	target transform.T

	Lines *lines.Mesh
	X     *Arrow
	Y     *Arrow
	Z     *Arrow

	XY *plane.Plane
	XZ *plane.Plane
	YZ *plane.Plane

	// screen size scaling factor
	size float32

	eye         vec3.T
	axis        vec3.T
	screenAxis  vec2.T
	start       vec2.T
	viewport    render.Screen
	vp          mat4.T
	proj        mat4.T
	scale       float32
	sensitivity float32
	dragging    bool
	hoverScale  vec3.T
}

var _ Gizmo = &Mover{}

// NewMover creates a new mover gizmo
func NewMover() *Mover {
	side := float32(0.2)
	planeAlpha := float32(0.3)

	s := side / 2

	mat := &material.Def{
		Shader:       "color_f",
		Subpass:      "forward",
		VertexFormat: vertex.C{},
		DepthTest:    true,
		DepthWrite:   true,
	}

	g := object.New("Mover Gizmo", &Mover{
		size:        0.12,
		sensitivity: 6,
		hoverScale:  vec3.New(1.1, 1.1, 1.1),

		X: object.Builder(NewArrow(color.Red)).
			Rotation(quat.Euler(0, 0, 270)).
			Create(),

		Y: NewArrow(color.Green),

		Z: object.Builder(NewArrow(color.Blue)).
			Rotation(quat.Euler(90, 0, 0)).
			Create(),

		// XY Plane
		XY: object.Builder(plane.NewObject(plane.Args{
			Mat:   mat,
			Size:  side,
			Color: color.Blue.WithAlpha(planeAlpha),
		})).
			Position(vec3.New(s, s, 0)).
			Rotation(quat.Euler(90, 0, 0)).
			Create(),

		// XZ Plane
		XZ: object.Builder(plane.NewObject(plane.Args{
			Mat:   mat,
			Size:  side,
			Color: color.Green.WithAlpha(planeAlpha),
		})).
			Rotation(quat.Euler(0, 90, 0)).
			Position(vec3.New(s, 0, s)).
			Create(),

		// YZ Plane
		YZ: object.Builder(plane.NewObject(plane.Args{
			Mat:   mat,
			Size:  side,
			Color: color.Red.WithAlpha(planeAlpha),
		})).
			Position(vec3.New(0, s, s)).
			Rotation(quat.Euler(0, 0, 90)).
			Create(),

		// Lines
		Lines: lines.New(lines.Args{
			Mat: &material.Def{
				Shader:       "lines",
				Subpass:      "output",
				VertexFormat: vertex.C{},
				Primitive:    vertex.Lines,
				DepthTest:    false,
			},
			Lines: []lines.Line{
				// xz lines
				lines.L(vec3.New(side, 0, 0), vec3.New(side, 0, side), color.Green),
				lines.L(vec3.New(side, 0, side), vec3.New(0, 0, side), color.Green),

				// xy lines
				lines.L(vec3.New(0, side, 0), vec3.New(side, side, 0), color.Blue),
				lines.L(vec3.New(side, 0, 0), vec3.New(side, side, 0), color.Blue),

				// yz lines
				lines.L(vec3.New(0, side, 0), vec3.New(0, side, side), color.Red),
				lines.L(vec3.New(0, 0, side), vec3.New(0, side, side), color.Red),
			},
		}),
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
		g.DragEnd(mouse.NopEvent())
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

func (g *Mover) DragStart(e mouse.Event, shape physics.Shape) {
	g.dragging = true

	g.axis = g.getColliderAxis(shape)
	cursor := g.viewport.NormalizeCursor(e.Position())
	g.start = cursor

	localDir := g.Transform().ProjectDir(g.axis)
	g.screenAxis = g.vp.TransformDir(localDir).XY().Normalized()
}

func (g *Mover) DragEnd(e mouse.Event) {
	g.dragging = false
}

func (g *Mover) DragMove(e mouse.Event) {
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

func (g *Mover) PreDraw(args render.Args, scene object.Object) error {
	g.eye = args.Position
	g.proj = args.Projection
	g.vp = args.VP
	g.viewport = args.Viewport
	return nil
}

func (g *Mover) Update(scene object.Component, dt float32) {
	g.Object.Update(scene, dt)

	// the gizmo should be displayed at the same size irrespectively of its distance to the camera.
	// we can undo the effects of perspective projection by measuring how much a vector would be "squeezed"
	// at the current distance form the camera, and then applying a scaling factor to counteract it.

	dist := vec3.Distance(g.Transform().WorldPosition(), g.eye)
	squeeze := g.proj.TransformPoint(vec3.New(1, 0, dist))
	f := g.size / squeeze.X
	g.scale = f
	g.Transform().SetWorldScale(vec3.New(f, f, f))
	g.Transform().SetWorldRotation(quat.Ident())
}

func (g *Mover) Dragging() bool          { return g.dragging }
func (g *Mover) Viewport() render.Screen { return g.viewport }
func (g *Mover) Camera() mat4.T          { return g.vp }

func (m *Mover) ToolMouseEvent(e mouse.Event, hover physics.RaycastHit) {
	HandleMouse(m, e, hover)
}
