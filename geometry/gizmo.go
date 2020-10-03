package geometry

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type Gizmo struct {
	*engine.Object
	X *engine.Object
	Y *engine.Object
	Z *engine.Object
}

func NewGizmo(position vec3.T) *Gizmo {
	radius := float32(0.1)
	height := float32(0.25)
	side := float32(0.2)
	segments := 6
	planeAlpha := float32(0.1)

	s := side / 2

	x := engine.NewObject(vec3.UnitX)
	x.Attach(NewCone(x, radius, height, segments, render.Red))
	x.Transform.Rotation.Z = -90

	xy := engine.NewObject(vec3.New(-s, s, 0))
	xy.Transform.Rotation.X = 90
	xy.Attach(NewPlane(xy, side, render.Blue.WithAlpha(planeAlpha)))

	y := engine.NewObject(vec3.UnitY)
	y.Attach(NewCone(y, radius, height, segments, render.Green))

	xz := engine.NewObject(vec3.New(-s, 0, -s))
	xz.Transform.Rotation.Y = 90
	xz.Attach(NewPlane(xz, side, render.Green.WithAlpha(planeAlpha)))

	z := engine.NewObject(vec3.UnitZ)
	z.Attach(NewCone(z, radius, height, segments, render.Blue))
	z.Transform.Rotation.X = 90

	yz := engine.NewObject(vec3.New(0, s, -s))
	yz.Transform.Rotation.Z = 90
	yz.Attach(NewPlane(yz, side, render.Red.WithAlpha(planeAlpha)))

	g := &Gizmo{
		Object: engine.NewObject(position),
		X:      x,
		Y:      y,
		Z:      z,
	}
	g.Attach(x, xy, xz, y, yz, z)

	lines := CreateLines(g.Object)

	// axis lines
	lines.Line(vec3.Zero, vec3.UnitX, render.Red)
	lines.Line(vec3.Zero, vec3.UnitY, render.Green)
	lines.Line(vec3.Zero, vec3.UnitZ, render.Blue)

	// xz lines
	lines.Line(vec3.Zero, vec3.New(-side, 0, 0), render.Green)
	lines.Line(vec3.New(-side, 0, 0), vec3.New(-side, 0, -side), render.Green)
	lines.Line(vec3.New(-side, 0, -side), vec3.New(0, 0, -side), render.Green)
	lines.Line(vec3.Zero, vec3.New(0, 0, -side), render.Green)

	// xy lines
	lines.Line(vec3.New(0, side, 0), vec3.New(-side, side, 0), render.Blue)
	lines.Line(vec3.New(-side, 0, 0), vec3.New(-side, side, 0), render.Blue)

	// yz lines
	lines.Line(vec3.New(0, side, 0), vec3.New(0, side, -side), render.Red)
	lines.Line(vec3.New(0, 0, -side), vec3.New(0, side, -side), render.Red)

	lines.Compute()
	g.Attach(lines)

	return g
}

func (g *Gizmo) Draw(args render.DrawArgs) {
	render.DepthMask(false)
	render.DepthTest(false)
	render.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	g.Object.Draw(args)

	render.DepthTest(true)
	render.DepthMask(true)
}
