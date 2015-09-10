package geometry

import (
    "unsafe"
    "github.com/go-gl/gl/v4.1-core/gl"
    mgl "github.com/go-gl/mathgl/mgl32"
    "github.com/johanhenriksson/goworld/render"
)

type LineVertex struct {
    X, Y, Z float32
    R, G, B, A float32
}

type LineVertices []LineVertex

func (buffer LineVertices) Elements() int { return len(buffer) }
func (buffer LineVertices) Size()     int { return 28 }
func (buffer LineVertices) GLPtr()    unsafe.Pointer { return gl.Ptr(buffer) }

type Lines struct {
    Lines       []Line
    Material    *render.Material
    Width       float32
    vao         *VertexArray
    vbo         *VertexBuffer
}

func CreateLines(mat *render.Material) *Lines {
    return &Lines {
        Lines:      make([]Line, 0),
        Material:   mat,
        Width:      10.0,
        vao:        CreateVertexArray(),
        vbo:        CreateVertexBuffer(),
    }
}

func (lines *Lines) Add(line Line) {
    lines.Lines = append(lines.Lines, line)
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
    data := make(LineVertices, 2 * count)
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
    lines.vbo.Buffer(data)
    lines.Material.Setup()
}

func (lines *Lines) Render() {
    gl.LineWidth(lines.Width)
    lines.Material.Use()
    lines.vao.Draw()
}

type Line struct {
    Start   mgl.Vec3
    End     mgl.Vec3
    Color   mgl.Vec4
}

func (lines *Lines) Line(start_x, start_y, start_z, end_x, end_y, end_z, r, g, b, a float32) {
    lines.Add(Line {
        Start: mgl.Vec3 { start_x, start_y, start_z },
        End:   mgl.Vec3 { end_x, end_y, end_z },
        Color: mgl.Vec4 { r, g, b, a },
    })
}

type Vec3 struct {
    X, Y, Z float32
}
