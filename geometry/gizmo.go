package geometry

import (
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// Gizmo is the visual representation of the 3D positioning tool
type Gizmo struct {
	*engine.Group
	Lines *Lines

	X  *Cone
	Y  *Cone
	Z  *Cone
	XY *Plane
	XZ *Plane
	YZ *Plane
}

// NewGizmo creates a new gizmo at the given position
func NewGizmo(position vec3.T) *Gizmo {
	radius := float32(0.1)
	height := float32(0.25)
	side := float32(0.2)
	segments := 6
	planeAlpha := float32(0.1)

	s := side / 2

	x := NewCone(radius, height, segments, render.Red)
	x.Position = vec3.UnitX
	x.Rotation = vec3.New(0, 0, -90)

	xy := NewPlane(side, render.Blue.WithAlpha(planeAlpha))
	xy.Position = vec3.New(-s, s, 0)
	xy.Rotation = vec3.New(90, 0, 0)

	y := NewCone(radius, height, segments, render.Green)
	y.Position = vec3.UnitY

	xz := NewPlane(side, render.Green.WithAlpha(planeAlpha))
	xz.Position = vec3.New(-s, 0, -s)
	xz.Rotation = vec3.New(0, 90, 0)

	z := NewCone(radius, height, segments, render.Blue)
	z.Position = vec3.UnitZ
	z.Rotation = vec3.New(90, 0, 0)

	yz := NewPlane(side, render.Red.WithAlpha(planeAlpha))
	yz.Position = vec3.New(0, s, -s)
	yz.Rotation = vec3.New(0, 0, 90)

	g := &Gizmo{
		Group: engine.NewGroup(position, vec3.Zero),
		Lines: CreateLines(),

		X:  x,
		Y:  y,
		Z:  z,
		XY: xy,
		XZ: xz,
		YZ: yz,
	}

	// axis lines
	g.Lines.Line(vec3.Zero, vec3.UnitX, render.Red)
	g.Lines.Line(vec3.Zero, vec3.UnitY, render.Green)
	g.Lines.Line(vec3.Zero, vec3.UnitZ, render.Blue)

	// xz lines
	g.Lines.Line(vec3.Zero, vec3.New(-side, 0, 0), render.Green)
	g.Lines.Line(vec3.New(-side, 0, 0), vec3.New(-side, 0, -side), render.Green)
	g.Lines.Line(vec3.New(-side, 0, -side), vec3.New(0, 0, -side), render.Green)
	g.Lines.Line(vec3.Zero, vec3.New(0, 0, -side), render.Green)

	// xy lines
	g.Lines.Line(vec3.New(0, side, 0), vec3.New(-side, side, 0), render.Blue)
	g.Lines.Line(vec3.New(-side, 0, 0), vec3.New(-side, side, 0), render.Blue)

	// yz lines
	g.Lines.Line(vec3.New(0, side, 0), vec3.New(0, side, -side), render.Red)
	g.Lines.Line(vec3.New(0, 0, -side), vec3.New(0, side, -side), render.Red)

	g.Lines.Compute()

	// attach components
	g.Attach(g.Lines, x, y, z, xy, xz, yz)

	return g
}
