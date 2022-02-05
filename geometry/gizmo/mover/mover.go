package mover

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/geometry/cone"
	"github.com/johanhenriksson/goworld/geometry/lines"
	"github.com/johanhenriksson/goworld/geometry/plane"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render/color"
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
		Color:    color.Red,
	}).
		Name("X Cone").
		Parent(g).
		Position(vec3.UnitX).
		Rotation(vec3.New(0, 0, 270)).
		Create()

	// Y arrow
	cone.Builder(&g.Y, cone.Args{
		Radius:   radius,
		Height:   height,
		Segments: segments,
		Color:    color.Green,
	}).
		Name("Y Cone").
		Parent(g).
		Position(vec3.UnitY).
		Create()

	// Z arrow
	cone.Builder(&g.Z, cone.Args{
		Radius:   radius,
		Height:   height,
		Segments: segments,
		Color:    color.Blue,
	}).
		Name("Z Cone").
		Parent(g).
		Position(vec3.UnitZN).
		Rotation(vec3.New(90, 180, 0)).
		Create()

	// XY plane
	plane.Builder(&g.XY, plane.Args{
		Size:  side,
		Color: color.Blue.WithAlpha(planeAlpha),
	}).
		Parent(g).
		Position(vec3.New(s, s, 0)).
		Rotation(vec3.New(90, 0, 0)).
		Create()

	// XZ plane
	plane.Builder(&g.XZ, plane.Args{
		Size:  side,
		Color: color.Green.WithAlpha(planeAlpha),
	}).
		Parent(g).
		Rotation(vec3.New(0, 90, 0)).
		Position(vec3.New(s, 0, s)).
		Create()

	// YZ plane
	plane.Builder(&g.YZ, plane.Args{
		Size:  side,
		Color: color.Red.WithAlpha(planeAlpha),
	}).
		Parent(g).
		Position(vec3.New(0, s, s)).
		Rotation(vec3.New(0, 0, 90)).
		Create()

	lines.Builder(&g.Lines, lines.Args{
		Lines: []lines.Line{
			// axis lines
			lines.L(vec3.Zero, vec3.UnitX, color.Red),
			lines.L(vec3.Zero, vec3.UnitY, color.Green),
			lines.L(vec3.Zero, vec3.UnitZN, color.Blue),

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
	}).
		Parent(g).
		Create()

	return g
}

func Attach(parent object.T, args Args) *T {
	box := New(args)
	parent.Adopt(box)
	return box
}

func NewObject(args Args) *T {
	parent := object.New("MoveGizmo")
	return Attach(parent, args)
}

func Builder(out **T, args Args) *object.Builder {
	b := object.Build("MoveGizmo")
	*out = New(args)
	return b.Adopt(*out)
}
