package main

import (
    "github.com/go-gl/gl/v4.1-core/gl"
//	"github.com/go-gl/mathgl/mgl32"

    "github.com/johanhenriksson/goworld/window"
    "github.com/johanhenriksson/goworld/shaders"
    "github.com/johanhenriksson/goworld/geometry"
)

func main() {
    wnd := window.Create("Hello World", 800, 600)

    program := shaders.CreateProgram()
    program.Attach(shaders.VertexShader("assets/shaders/trivial.vs.glsl"))
    program.Attach(shaders.FragmentShader("assets/shaders/trivial.fs.glsl"))
    program.Link()

    vtx := []geometry.Vertex {
        geometry.Vertex { X:  0.0, Y:  0.5, Z: 0, R: 1.0, G: 0.0, B: 0.0 },
        geometry.Vertex { X:  0.5, Y: -0.5, Z: 0, R: 0.0, G: 1.0, B: 0.0 },
        geometry.Vertex { X: -0.5, Y: -0.5, Z: 0, R: 0.0, G: 0.0, B: 1.0 },
    }

    vao := geometry.CreateVertexArray()
    vao.Bind()

    vbo := geometry.CreateVertexBuffer()
    vbo.Bind()
    vbo.Buffer(vtx)

    vertexAttr := program.GetAttributeLocation("vertex")
	gl.EnableVertexAttribArray(vertexAttr)
    gl.VertexAttribPointer(vertexAttr, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))

    colorAttr  := program.GetAttributeLocation("color")
	gl.EnableVertexAttribArray(colorAttr)
	gl.VertexAttribPointer(colorAttr, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))

    /*
	proj   := mgl32.Perspective(mgl32.DegToRad(45.0), float32(800.0/600.0), 0.1, 10.0)
	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	model  := mgl32.Ident4()
    */

    //gl.ClearColor(1.0, 0.0, 0.0, 1.0)
    wnd.SetRenderCallback(func(wnd *window.Window, dt float32) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        program.Use()
        vao.Bind()

		gl.DrawArrays(gl.TRIANGLES, 0, 3)
    })

    wnd.Loop()
}
