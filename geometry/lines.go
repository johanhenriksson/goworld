package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

type Lines struct {
	// *engine.Object
	Lines    []Line
	Material *render.Material
	vao      *render.VertexArray
	name     string
}

type Line struct {
	Start vec3.T
	End   vec3.T
	Color render.Color
}

func CreateLines(name string, lines ...Line) *Lines {
	l := &Lines{
		// Object:   parent,
		Lines:    lines,
		Material: assets.GetMaterialShared("lines"),
		vao:      render.CreateVertexArray(render.Lines, "geometry"),
		name:     name,
	}
	l.Compute()
	return l
}

func (lines *Lines) Name() string {
	return lines.name
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
		a.Position = line.Start
		a.Color = line.Color
		b.Position = line.End
		b.Color = line.Color
	}

	ptr := lines.Material.VertexPointers(data)
	lines.vao.BufferTo(ptr, data)
}

func (lines *Lines) Update(dt float32) {}

func (lines *Lines) DrawLines(args engine.DrawArgs) {
	// setup line material
	if len(lines.Lines) > 0 && args.Pass == render.LinePass {
		lines.Material.Use()
		lines.Material.Mat4("mvp", &args.MVP)
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

func (lines *Lines) Collect(pass engine.DrawPass, args engine.DrawArgs) {
	if pass.Visible(lines, args) {
		pass.Queue(lines, args)
	}
}
