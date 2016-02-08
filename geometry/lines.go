package geometry

import (
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/assets"

    "github.com/go-gl/gl/v4.1-core/gl"
    mgl "github.com/go-gl/mathgl/mgl32"
)

type Lines struct {
    Lines       []Line
    Material    *render.Material
    Width       float32
    vao         *render.VertexArray
    vbo         *render.VertexBuffer
}

type Line struct {
    Start   mgl.Vec3
    End     mgl.Vec3
    Color   mgl.Vec4
}

func CreateLines() *Lines {
    l := &Lines {
        Lines:    make([]Line, 0, 0),
        Material: assets.GetMaterialCached("lines"),
        Width:    10.0,
        vao:      render.CreateVertexArray(),
        vbo:      render.CreateVertexBuffer(),
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

func (lines *Lines) Box(x,y,z,w,h,d,r,g,b,a float32) {
    /* Bottom square */
    lines.Line(x,y,z, x+w, y, z, r,g,b,a)
    lines.Line(x,y,z, x, y, z+d, r,g,b,a)
    lines.Line(x+w,y,z,x+w,y,z+d, r,g,b,a)
    lines.Line(x,y,z+w,x+w,y,z+d, r,g,b,a)

    /* Top square */
    lines.Line(x,y+h,z, x+w, y+h, z, r,g,b,a)
    lines.Line(x,y+h,z, x, y+h, z+d, r,g,b,a)
    lines.Line(x+w,y+h,z,x+w,y+h,z+d, r,g,b,a)
    lines.Line(x,y+h,z+w,x+w,y+h,z+d, r,g,b,a)

    lines.Line(x,y,z, x, y+h, z, r,g,b,a)
    lines.Line(x+w,y,z, x+w, y+h, z, r,g,b,a)
    lines.Line(x,y,z+d, x, y+h, z+d, r,g,b,a)
    lines.Line(x+w,y,z+d, x+w, y+h, z+d, r,g,b,a)
}

func (lines *Lines) Compute() {
    count := len(lines.Lines)
    data := make(ColorVertices, 2 * count)
    for i := 0; i < count; i++ {
        line := lines.Lines[i]
        a := &data[2*i+0]
        b := &data[2*i+1]
        a.X, a.Y, a.Z      = line.Start[0], line.Start[1], line.Start[2]
        b.X, b.Y, b.Z      = line.End[0],   line.End[1],   line.End[2]
        a.R, a.G, a.B, a.A = line.Color[0], line.Color[1], line.Color[2], line.Color[3]
        b.R, b.G, b.B, b.A = line.Color[0], line.Color[1], line.Color[2], line.Color[3]
    }
    lines.vao.Length = int32(2 * count)
    lines.vao.Type   = gl.LINES
    lines.vao.Bind()
    lines.vbo.Bind()
    if lines.vao.Length > 0 {
        lines.vbo.Buffer(data)
    }
    lines.Material.SetupVertexPointers()
}

func (lines *Lines) Draw(args render.DrawArgs) {
    // setup line material
    if len(lines.Lines) > 0 && args.Pass == "lines" {
        lines.vao.Draw()
    }
}

func (lines *Lines) Line(start_x, start_y, start_z, end_x, end_y, end_z, r, g, b, a float32) {
    lines.Add(Line {
        Start: mgl.Vec3 { start_x, start_y, start_z },
        End:   mgl.Vec3 { end_x, end_y, end_z },
        Color: mgl.Vec4 { r, g, b, a },
    })
}
