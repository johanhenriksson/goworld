package gizmo

import (
	"github.com/johanhenriksson/goworld/core/collider"
	"github.com/johanhenriksson/goworld/core/input/mouse"
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/transform"
	"github.com/johanhenriksson/goworld/geometry/cube"
	"github.com/johanhenriksson/goworld/geometry/cyllinder"
	"github.com/johanhenriksson/goworld/math/mat4"
	"github.com/johanhenriksson/goworld/math/vec2"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/vertex"
)

// Scaler Gizmo can be used to scale objects in the 3D scene.
type Scaler struct {
	object.T

	target transform.T

	All *cube.T
	X   *cube.T
	Xb  *cyllinder.T
	Y   *cube.T
	Yb  *cyllinder.T
	Z   *cube.T
	Zb  *cyllinder.T

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

var _ Gizmo = &Scaler{}

// NewScaled creates a new scaler gizmo
func NewScaler() *Scaler {
	size := float32(0.2)
	bodyRadius := float32(0.025)
	segments := 16

	mat := &material.Def{
		Shader:       "color_f",
		Subpass:      "forward",
		VertexFormat: vertex.C{},
		DepthTest:    true,
		DepthWrite:   true,
	}

	g := object.New(&Scaler{
		size:        0.12,
		sensitivity: 6,
		hoverScale:  vec3.New(1.2, 1.2, 1.2),

		// All Axis Cube
		All: object.Builder(cube.New(cube.Args{
			Mat:   mat,
			Size:  size,
			Color: color.RGB(0.7, 0.7, 0.7),
		})).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*size, 2*size, 2*size),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// X Arrow Cube
		X: object.Builder(cube.New(cube.Args{
			Mat:   mat,
			Size:  size,
			Color: color.Red,
		})).
			Position(vec3.UnitX).
			Rotation(vec3.New(0, 0, 270)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*size, 2*size, 2*size),
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
				Size:   vec3.New(size, 1, size),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// Y Arrow Cube
		Y: object.Builder(cube.New(cube.Args{
			Mat:   mat,
			Size:  size,
			Color: color.Green,
		})).
			Position(vec3.UnitY).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*size, 2*size, 2*size),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// Y Arrow Body
		Yb: object.Builder(cyllinder.New(cyllinder.Args{
			Mat:      mat,
			Radius:   bodyRadius,
			Height:   1,
			Segments: segments,
			Color:    color.Green,
		})).
			Position(vec3.New(0, 0.5, 0)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*size, 1, 2*size),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),

		// Z Arrow Cube
		Z: object.Builder(cube.New(cube.Args{
			Mat:   mat,
			Size:  size,
			Color: color.Blue,
		})).
			Position(vec3.UnitZ).
			Rotation(vec3.New(90, 180, 0)).
			Attach(collider.NewBox(collider.Box{
				Size:   vec3.New(2*size, 2*size, 2*size),
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
				Size:   vec3.New(2*size, 1, 2*size),
				Center: vec3.New(0, 0, 0),
			})).
			Create(),
	})

	return g
}

func (g *Scaler) Name() string {
	return "ScalerGizmo"
}

func (g *Scaler) Target() transform.T {
	return g.target
}

func (g *Scaler) SetTarget(t transform.T) {
	if t != nil {
		g.Transform().SetPosition(t.WorldPosition())
	}
	g.target = t
}

func (g *Scaler) CanDeselect() bool {
	return true
}

func (g *Scaler) getColliderAxis(collider collider.T) vec3.T {
	axisObj := collider.Parent()
	switch axisObj {
	case g.All:
		return vec3.One
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

func (g *Scaler) DragStart(e mouse.Event, collider collider.T) {
	g.dragging = true

	g.axis = g.getColliderAxis(collider)
	cursor := g.viewport.NormalizeCursor(e.Position())
	g.start = cursor

	localDir := g.Transform().ProjectDir(g.axis)
	g.screenAxis = g.vp.TransformDir(localDir).XY().Normalized()
}

func (g *Scaler) DragEnd(e mouse.Event) {
	g.dragging = false
}

func (g *Scaler) DragMove(e mouse.Event) {
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

func (g *Scaler) Hover(hovering bool, collider collider.T) {
	if hovering {
		// hover start
		axis := g.getColliderAxis(collider)
		switch axis {
		case vec3.One:
			g.All.Transform().SetScale(g.hoverScale)
			g.X.Transform().SetScale(g.hoverScale)
			g.Y.Transform().SetScale(g.hoverScale)
			g.Z.Transform().SetScale(g.hoverScale)
		case vec3.UnitX:
			g.All.Transform().SetScale(vec3.One)
			g.X.Transform().SetScale(g.hoverScale)
			g.Y.Transform().SetScale(vec3.One)
			g.Z.Transform().SetScale(vec3.One)
		case vec3.UnitY:
			g.All.Transform().SetScale(vec3.One)
			g.X.Transform().SetScale(vec3.One)
			g.Y.Transform().SetScale(g.hoverScale)
			g.Z.Transform().SetScale(vec3.One)
		case vec3.UnitZ:
			g.All.Transform().SetScale(vec3.One)
			g.X.Transform().SetScale(vec3.One)
			g.Y.Transform().SetScale(vec3.One)
			g.Z.Transform().SetScale(g.hoverScale)
		}
	}
	if !hovering {
		// reset scaling
		g.All.Transform().SetScale(vec3.One)
		g.X.Transform().SetScale(vec3.One)
		g.Y.Transform().SetScale(vec3.One)
		g.Z.Transform().SetScale(vec3.One)
	}
}

func (g *Scaler) PreDraw(args render.Args, scene object.T) error {
	g.eye = args.Position
	g.proj = args.Projection
	g.vp = args.VP
	g.viewport = args.Viewport
	return nil
}

func (g *Scaler) Update(scene object.T, dt float32) {
	g.T.Update(scene, dt)

	// the gizmo should be displayed at the same size irrespectively of its distance to the camera.
	// we can undo the effects of perspective projection by measuring how much a vector would be "squeezed"
	// at the current distance form the camera, and then applying a scaling factor to counteract it.

	dist := vec3.Distance(g.Transform().WorldPosition(), g.eye)
	squeeze := g.proj.TransformPoint(vec3.New(1, 0, dist))
	f := g.size / squeeze.X
	g.scale = f
	g.Transform().SetScale(vec3.New(f, f, f))
}

func (g *Scaler) Dragging() bool          { return g.dragging }
func (g *Scaler) Viewport() render.Screen { return g.viewport }
func (g *Scaler) Camera() mat4.T          { return g.vp }

func (m *Scaler) MouseEvent(e mouse.Event) {
	HandleMouse(m, e)
}
