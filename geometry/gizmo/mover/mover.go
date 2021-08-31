package mover

import (
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/geometry/cone"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Mover Gizmo is the visual representation of the 3D positioning tool
type T struct {
	object.T
	Args

	Lines *lines.T
	X     *cone.T
	Y     *cone.T
	Z     *cone.T
	XY    *plane.T
	XZ    *plane.T
	YZ    *plane.T
}

type Args struct {
}

// New creates a new gizmo at the given position
func New(args Args) *T {
	radius := float32(0.1)
	height := float32(0.25)
	side := float32(0.2)
	segments := 6
	planeAlpha := float32(0.3)

	s := side / 2

	g := &T{
		T: object.New("Gizmo"),
	}

	// X arrow
	cone.Builder(&g.X, cone.Args{
		Radius:   radius,
		Height:   height,
		Segments: segments,
		Color:    render.Red,
	}).
		Position(vec3.UnitX).
		Rotation(vec3.New(0, 0, -90)).
		Create(g.T)

	// Y arrow
	cone.Builder(&g.Y, cone.Args{
		Radius:   radius,
		Height:   height,
		Segments: segments,
		Color:    render.Green,
	}).
		Position(vec3.New(0, 1, 0)).
		Create(g.T)

	// Z arrow
	cone.Builder(&g.Z, cone.Args{
		Radius:   radius,
		Height:   height,
		Segments: segments,
		Color:    render.Blue,
	}).
		Position(vec3.UnitZ).
		Rotation(vec3.New(90, 0, 0)).
		Create(g.T)

	// XY plane
	plane.Builder(&g.XY, plane.Args{
		Size:  side,
		Color: render.Blue.WithAlpha(planeAlpha),
	}).
		Position(vec3.New(s, s, 0)).
		Rotation(vec3.New(90, 0, 0)).
		Create(g.T)

	// XZ plane
	plane.Builder(&g.XZ, plane.Args{
		Size:  side,
		Color: render.Green.WithAlpha(planeAlpha),
	}).
		Rotation(vec3.New(0, 90, 0)).
		Position(vec3.New(s, 0, s)).
		Create(g.T)

	// YZ plane
	plane.Builder(&g.YZ, plane.Args{
		Size:  side,
		Color: render.Red.WithAlpha(planeAlpha),
	}).
		Position(vec3.New(0, s, s)).
		Rotation(vec3.New(0, 0, 90)).
		Create(g.T)

	lines.Builder(&g.Lines, lines.Args{
		Lines: []lines.Line{
			// axis lines
			lines.L(vec3.Zero, vec3.UnitX, render.Red),
			lines.L(vec3.Zero, vec3.UnitY, render.Green),
			lines.L(vec3.Zero, vec3.UnitZ, render.Blue),

			// xz lines
			lines.L(vec3.New(side, 0, 0), vec3.New(side, 0, side), render.Green),
			lines.L(vec3.New(side, 0, side), vec3.New(0, 0, side), render.Green),

			// xy lines
			lines.L(vec3.New(0, side, 0), vec3.New(side, side, 0), render.Blue),
			lines.L(vec3.New(side, 0, 0), vec3.New(side, side, 0), render.Blue),

			// yz lines
			lines.L(vec3.New(0, side, 0), vec3.New(0, side, side), render.Red),
			lines.L(vec3.New(0, 0, side), vec3.New(0, side, side), render.Red),
		},
	}).Create(g.T)

	return g
}

func Attach(parent object.T, args Args) *T {
	box := New(args)
	parent.Attach(box)
	return box
}

func NewObject(args Args) *T {
	parent := object.New("MoveGizmo")
	return Attach(parent, args)
}

func Builder(out **T, args Args) *object.Builder {
	b := object.Build("MoveGizmo")
	*out = New(args)
	return b.Attach(*out)
}
