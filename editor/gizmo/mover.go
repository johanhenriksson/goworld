package gizmo

import (
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/geometry/cone"
	"github.com/johanhenriksson/goworld/geometry/cyllinder"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Mover Gizmo can be used to reposition objects in the 3D scene.
type Mover struct {
	object.T

	target transform.T

	Lines *lines.T
	X     *cone.T
	Xb    *cyllinder.T
	Y     *cone.T
	Yb    *cyllinder.T
	Z     *cone.T
	Zb    *cyllinder.T
	XY    *plane.T
	XZ    *plane.T
	YZ    *plane.T

	axis       vec3.T
	screenAxis vec2.T
	start      vec2.T
	viewport   render.Screen
	camera     mat4.T
	dragging   bool
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
		// X Arrow Cone
		X: object.Builder(cone.New(cone.Args{
			Mat:      mat,
			Radius:   radius,
			Height:   height,
			Segments: segments,
			Color:    color.Red,
		})).
			Position(vec3.UnitX).
			Rotation(vec3.New(0, 0, 270)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*radius, height, 2*radius),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// X Arrow Body
		Xb: object.Builder(cyllinder.New(cyllinder.Args{
			Mat:      mat,
			Radius:   bodyRadius,
			Height:   1,
			Segments: segments,
			Color:    color.Red,
		})).
			Position(vec3.New(0.5, 0, 0)).
			Rotation(vec3.New(0, 0, 270)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(radius, 1, radius),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// Y Arrow Cone
		Y: object.Builder(cone.New(cone.Args{
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
		Yb: object.Builder(cyllinder.New(cyllinder.Args{
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
		Z: object.Builder(cone.New(cone.Args{
			Mat:      mat,
			Radius:   radius,
			Height:   height,
			Segments: segments,
			Color:    color.Blue,
		})).
			Position(vec3.UnitZ).
			Rotation(vec3.New(90, 180, 0)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*radius, height, 2*radius),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// Z Arrow Body
		Zb: object.Builder(cyllinder.New(cyllinder.Args{
			Mat:      mat,
			Radius:   bodyRadius,
			Height:   1,
			Segments: segments,
			Color:    color.Blue,
		})).
			Position(vec3.New(0, 0, 0.5)).
			Rotation(vec3.New(90, 180, 0)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*radius, 1, 2*radius),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// XY Plane
		XY: object.Builder(plane.New(plane.Args{
			Mat:   mat,
			Size:  side,
			Color: color.Blue.WithAlpha(planeAlpha),
		})).
			Position(vec3.New(s, s, 0)).
			Rotation(vec3.New(90, 0, 0)).
			Create(),

		// XZ Plane
		XZ: object.Builder(plane.New(plane.Args{
			Mat:   mat,
			Size:  side,
			Color: color.Green.WithAlpha(planeAlpha),
		})).
			Rotation(vec3.New(0, 90, 0)).
			Position(vec3.New(s, 0, s)).
			Create(),

		// YZ Plane
		YZ: object.Builder(plane.New(plane.Args{
			Mat:   mat,
			Size:  side,
			Color: color.Red.WithAlpha(planeAlpha),
		})).
			Position(vec3.New(0, s, s)).
			Rotation(vec3.New(0, 0, 90)).
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

func (g *Mover) DragStart(e mouse.Event, collider collider.T) {
	g.dragging = true

	axisObj := collider.Parent()
	switch axisObj {
	case g.X:
		fallthrough
	case g.Xb:
		g.axis = vec3.UnitX
	case g.Y:
		fallthrough
	case g.Yb:
		g.axis = vec3.UnitY
	case g.Z:
		fallthrough
	case g.Zb:
		g.axis = vec3.UnitZ
	default:
		return
	}
	cursor := g.viewport.NormalizeCursor(e.Position())
	g.start = cursor

	localDir := g.Transform().ProjectDir(g.axis)
	g.screenAxis = g.camera.TransformDir(localDir).XY().Normalized()
}

func (g *Mover) DragEnd(e mouse.Event) {
	g.dragging = false
}

func (g *Mover) DragMove(e mouse.Event) {
	if e.Action() == mouse.Move {
		cursor := g.viewport.NormalizeCursor(e.Position())

		delta := g.start.Sub(cursor)
		mag := -5 * vec2.Dot(delta, g.screenAxis) / g.screenAxis.Length()
		g.start = cursor
		pos := g.Transform().Position().Add(g.axis.Scaled(mag))
		g.Transform().SetPosition(pos)

		if g.target != nil {
			g.target.SetWorldPosition(pos)
		}
	}
}

func (g *Mover) PreDraw(args render.Args, scene object.T) error {
	g.camera = args.VP
	g.viewport = args.Viewport
	return nil
}

func (g *Mover) Dragging() bool          { return g.dragging }
func (g *Mover) Viewport() render.Screen { return g.viewport }
func (g *Mover) Camera() mat4.T          { return g.camera }

func (m *Mover) MouseEvent(e mouse.Event) {
	HandleMouse(m, e)
}
