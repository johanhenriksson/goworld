package gizmo

import (
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/geometry/cone"
	"github.com/johanhenriksson/goworld/geometry/cylinder"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/quat"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Mover Gizmo can be used to reposition objects in the 3D scene.
type Mover struct {
	object.G

	target transform.T

	Lines *lines.T
	X     *cone.Cone
	Xb    *cylinder.Cylinder
	Y     *cone.Cone
	Yb    *cylinder.Cylinder
	Z     *cone.Cone
	Zb    *cylinder.Cylinder
	XY    *plane.Plane
	XZ    *plane.Plane
	YZ    *plane.Plane

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
	radius := float32(0.1)
	bodyRadius := radius / 4
	height := float32(0.35)
	side := float32(0.2)
	segments := 32
	planeAlpha := float32(0.3)

	s := side / 2

	mat := &material.Def{
		Shader:       "color_f",
		Subpass:      "forward",
		VertexFormat: vertex.C{},
		DepthTest:    true,
		DepthWrite:   true,
	}

	g := object.New(&Mover{
		size:        0.12,
		sensitivity: 6,
		hoverScale:  vec3.New(1.2, 1.2, 1.2),

		// X Arrow Cone
		X: object.Builder(
			cone.Group(cone.Args{
				Mat:      mat,
				Radius:   radius,
				Height:   height,
				Segments: segments,
				Color:    color.Red,
			})).
			Position(vec3.UnitX).
			Rotation(quat.Euler(0, 0, 270)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*radius, height, 2*radius),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// X Arrow Body
		Xb: object.Builder(cylinder.Group(cylinder.Args{
			Mat:      mat,
			Radius:   bodyRadius,
			Height:   1,
			Segments: segments,
			Color:    color.Red,
		})).
			Position(vec3.New(0.5, 0, 0)).
			Rotation(quat.Euler(0, 0, 270)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(radius, 1, radius),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// Y Arrow Cone
		Y: object.Builder(cone.Group(cone.Args{
			Mat:      mat,
			Radius:   radius,
			Height:   height,
			Segments: segments,
			Color:    color.Green,
		})).
			Position(vec3.UnitY).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*radius, height, 2*radius),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// Y Arrow body
		Yb: object.Builder(cylinder.Group(cylinder.Args{
			Mat:      mat,
			Radius:   bodyRadius,
			Height:   1,
			Segments: segments,
			Color:    color.Green,
		})).
			Position(vec3.New(0, 0.5, 0)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*radius, 1, 2*radius),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// Z Arrow Cone
		Z: object.Builder(cone.Group(cone.Args{
			Mat:      mat,
			Radius:   radius,
			Height:   height,
			Segments: segments,
			Color:    color.Blue,
		})).
			Position(vec3.UnitZ).
			Rotation(quat.Euler(90, 180, 0)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*radius, height, 2*radius),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// Z Arrow Body
		Zb: object.Builder(cylinder.Group(cylinder.Args{
			Mat:      mat,
			Radius:   bodyRadius,
			Height:   1,
			Segments: segments,
			Color:    color.Blue,
		})).
			Position(vec3.New(0, 0, 0.5)).
			Rotation(quat.Euler(90, 180, 0)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*radius, 1, 2*radius),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// XY Plane
		XY: object.Builder(plane.Group(plane.Args{
			Mat:   mat,
			Size:  side,
			Color: color.Blue.WithAlpha(planeAlpha),
		})).
			Position(vec3.New(s, s, 0)).
			Rotation(quat.Euler(90, 0, 0)).
			Create(),

		// XZ Plane
		XZ: object.Builder(plane.Group(plane.Args{
			Mat:   mat,
			Size:  side,
			Color: color.Green.WithAlpha(planeAlpha),
		})).
			Rotation(quat.Euler(0, 90, 0)).
			Position(vec3.New(s, 0, s)).
			Create(),

		// YZ Plane
		YZ: object.Builder(plane.Group(plane.Args{
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

func (g *Mover) Name() string {
	return "MoverGizmo"
}

func (g *Mover) Target() transform.T {
	return g.target
}

func (g *Mover) SetTarget(t transform.T) {
	if t != nil {
		g.Transform().SetPosition(t.WorldPosition())
	}
	g.target = t
}

func (g *Mover) CanDeselect() bool {
	return true
}

func (g *Mover) getColliderAxis(collider collider.T) vec3.T {
	axisObj := collider.Parent()
	switch axisObj {
	case g.X:
		fallthrough
	case g.Xb:
		return vec3.UnitX
	case g.Y:
		fallthrough
	case g.Yb:
		return vec3.UnitY
	case g.Z:
		fallthrough
	case g.Zb:
		return vec3.UnitZ
	}
	return vec3.Zero
}

func (g *Mover) DragStart(e mouse.Event, collider collider.T) {
	g.dragging = true

	g.axis = g.getColliderAxis(collider)
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

		delta := g.start.Sub(cursor)
		mag := -1 * g.sensitivity * g.scale * vec2.Dot(delta, g.screenAxis) / g.screenAxis.Length()
		g.start = cursor
		pos := g.Transform().Position().Add(g.axis.Scaled(mag))
		g.Transform().SetPosition(pos)

		if g.target != nil {
			g.target.SetWorldPosition(pos)
		}
	}
}

func (g *Mover) Hover(hovering bool, collider collider.T) {
	if hovering {
		// hover start
		axis := g.getColliderAxis(collider)
		switch axis {
		case vec3.UnitX:
			g.X.Transform().SetScale(g.hoverScale)
			g.Y.Transform().SetScale(vec3.One)
			g.Z.Transform().SetScale(vec3.One)
		case vec3.UnitY:
			g.X.Transform().SetScale(vec3.One)
			g.Y.Transform().SetScale(g.hoverScale)
			g.Z.Transform().SetScale(vec3.One)
		case vec3.UnitZ:
			g.X.Transform().SetScale(vec3.One)
			g.Y.Transform().SetScale(vec3.One)
			g.Z.Transform().SetScale(g.hoverScale)
		}
	}
	if !hovering {
		// reset scaling
		g.X.Transform().SetScale(vec3.One)
		g.Y.Transform().SetScale(vec3.One)
		g.Z.Transform().SetScale(vec3.One)
	}
}

func (g *Mover) PreDraw(args render.Args, scene object.T) error {
	g.eye = args.Position
	g.proj = args.Projection
	g.vp = args.VP
	g.viewport = args.Viewport
	return nil
}

func (g *Mover) Update(scene object.T, dt float32) {
	g.G.Update(scene, dt)

	// the gizmo should be displayed at the same size irrespectively of its distance to the camera.
	// we can undo the effects of perspective projection by measuring how much a vector would be "squeezed"
	// at the current distance form the camera, and then applying a scaling factor to counteract it.

	dist := vec3.Distance(g.Transform().WorldPosition(), g.eye)
	squeeze := g.proj.TransformPoint(vec3.New(1, 0, dist))
	f := g.size / squeeze.X
	g.scale = f
	g.Transform().SetScale(vec3.New(f, f, f))
}

func (g *Mover) Dragging() bool          { return g.dragging }
func (g *Mover) Viewport() render.Screen { return g.viewport }
func (g *Mover) Camera() mat4.T          { return g.vp }

func (m *Mover) MouseEvent(e mouse.Event) {
	HandleMouse(m, e)
}
