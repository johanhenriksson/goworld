package geometry

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/engine"
	"github.com/johanhenriksson/goworld/engine/object"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/vertex"
)

type Lines struct {
	*object.T
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

func NewLines(name string, lines ...Line) *Lines {
	l := &Lines{
		T:        object.New("Lines"),
		Lines:    lines,
		Material: assets.GetMaterialShared("lines"),
		vao:      render.CreateVertexArray(render.Lines),
		name:     name,
	}
	l.Compute()
	return l
}

// func (lines *Lines) Parent() *engine.Object     { return lines.Object }
// func (lines *Lines) SetParent(p *engine.Object) { lines.Object = p }
// func (lines *Lines) Collect(q *engine.Query)    {}

func (li *Lines) String() string {
	return li.name
}

func (li *Lines) Add(line Line) {
	li.Lines = append(li.Lines, line)
}

func (li *Lines) Clear() {
	li.Lines = make([]Line, 0, 0)
	li.Compute()
}

func (li *Lines) Box(x, y, z, w, h, d float32, color render.Color) {
	// bottom square
	li.Line(vec3.New(x, y, z), vec3.New(x+w, y, z), color)
	li.Line(vec3.New(x, y, z), vec3.New(x, y, z+d), color)
	li.Line(vec3.New(x+w, y, z), vec3.New(x+w, y, z+d), color)
	li.Line(vec3.New(x, y, z+w), vec3.New(x+w, y, z+d), color)

	// top square
	li.Line(vec3.New(x, y+h, z), vec3.New(x+w, y+h, z), color)
	li.Line(vec3.New(x, y+h, z), vec3.New(x, y+h, z+d), color)
	li.Line(vec3.New(x+w, y+h, z), vec3.New(x+w, y+h, z+d), color)
	li.Line(vec3.New(x, y+h, z+w), vec3.New(x+w, y+h, z+d), color)

	li.Line(vec3.New(x, y, z), vec3.New(x, y+h, z), color)
	li.Line(vec3.New(x+w, y, z), vec3.New(x+w, y+h, z), color)
	li.Line(vec3.New(x, y, z+d), vec3.New(x, y+h, z+d), color)
	li.Line(vec3.New(x+w, y, z+d), vec3.New(x+w, y+h, z+d), color)
}

func (li *Lines) Compute() {
	count := len(li.Lines)
	data := make([]vertex.C, 2*count)
	for i := 0; i < count; i++ {
		line := li.Lines[i]
		a := &data[2*i+0]
		b := &data[2*i+1]
		a.P = line.Start
		a.C = line.Color.Vec4()
		b.P = line.End
		b.C = line.Color.Vec4()
	}

	ptr := li.Material.VertexPointers(data)
	li.vao.BufferTo(ptr, data)
}

func (li *Lines) DrawLines(args engine.DrawArgs) {
	// setup line material
	if len(li.Lines) > 0 && args.Pass == render.Line {
		li.Material.Use()
		li.Material.Mat4("mvp", &args.MVP)
		li.vao.Draw()
	}
}

func (li *Lines) Line(start, end vec3.T, color render.Color) {
	li.Add(Line{
		Start: start,
		End:   end,
		Color: color,
	})
}
