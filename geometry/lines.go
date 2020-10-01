package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type Lines struct {
	Lines    []Line
	Material *render.Material
	vao      *render.VertexArray
}

type Line struct {
	Start vec3.T
	End   vec3.T
	Color render.Color
}

func CreateLines() *Lines {
	l := &Lines{
		Lines:    make([]Line, 0, 8),
		Material: assets.GetMaterialCached("lines"),
		vao:      render.CreateVertexArray(render.Lines, "geometry"),
	}
	l.Compute()
	return l
}

func (lines *Lines) Add(line Line) {
	lines.Lines = append(lines.Lines, line)
}

func (lines *Lines) Clear() {
	lines.Lines = make([]Line, 0, 0)
	lines.Compute()
}

func (lines *Lines) Box(x, y, z, w, h, d float32, color render.Color) {
	/* Bottom square */
	lines.Line(vec3.New(x, y, z), vec3.New(x+w, y, z), color)
	lines.Line(vec3.New(x, y, z), vec3.New(x, y, z+d), color)
	lines.Line(vec3.New(x+w, y, z), vec3.New(x+w, y, z+d), color)
	lines.Line(vec3.New(x, y, z+w), vec3.New(x+w, y, z+d), color)

	/* Top square */
	lines.Line(vec3.New(x, y+h, z), vec3.New(x+w, y+h, z), color)
	lines.Line(vec3.New(x, y+h, z), vec3.New(x, y+h, z+d), color)
	lines.Line(vec3.New(x+w, y+h, z), vec3.New(x+w, y+h, z+d), color)
	lines.Line(vec3.New(x, y+h, z+w), vec3.New(x+w, y+h, z+d), color)

	lines.Line(vec3.New(x, y, z), vec3.New(x, y+h, z), color)
	lines.Line(vec3.New(x+w, y, z), vec3.New(x+w, y+h, z), color)
	lines.Line(vec3.New(x, y, z+d), vec3.New(x, y+h, z+d), color)
	lines.Line(vec3.New(x+w, y, z+d), vec3.New(x+w, y+h, z+d), color)
}

func (lines *Lines) Compute() {
	count := len(lines.Lines)
	data := make(ColorVertices, 2*count)
	for i := 0; i < count; i++ {
		line := lines.Lines[i]
		a := &data[2*i+0]
		b := &data[2*i+1]
		a.X, a.Y, a.Z = line.Start.X, line.Start.Y, line.Start.Z
		b.X, b.Y, b.Z = line.End.X, line.End.Y, line.End.Z
		a.R, a.G, a.B, a.A = line.Color.R, line.Color.G, line.Color.B, line.Color.A
		b.R, b.G, b.B, b.A = line.Color.R, line.Color.G, line.Color.B, line.Color.A
	}
	lines.vao.Bind()
	if len(data) > 0 {
		lines.vao.Buffer("geometry", data)
	} else {
		lines.vao.Length = 0
	}
	lines.Material.SetupVertexPointers()
}

func (lines *Lines) Draw(args render.DrawArgs) {
	// setup line material
	if len(lines.Lines) > 0 && args.Pass == render.LinePass {
		lines.Material.Use()
		lines.Material.Mat4f("mvp", &args.MVP)
		lines.vao.Draw()
	}
}

func (lines *Lines) Line(start, end vec3.T, color render.Color) {
	lines.Add(Line{
		Start: start,
		End:   end,
		Color: color,
	})
}
