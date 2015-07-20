package main

import (
    "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

    "github.com/johanhenriksson/goworld/window"
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/geometry"
)

func main() {
    wnd := window.Create("Hello World", 800, 600)

    program := render.CreateProgram()
    program.Attach(render.VertexShader("assets/shaders/vertex.glsl"))
    program.Attach(render.FragmentShader("assets/shaders/fragment.glsl"))
    program.Link()
    program.Use()

    cube := geometry.Cube()
    cube.Bind()

    tx, _ := render.LoadTexture("assets/square.png")
    tx.Bind(0)

    gl.CullFace(gl.FRONT_AND_BACK)

    vertexAttr := program.GetAttributeLocation("vertex")
	gl.EnableVertexAttribArray(vertexAttr)
    gl.VertexAttribPointer(vertexAttr, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

    /*
    colorAttr  := program.GetAttributeLocation("color")
	gl.EnableVertexAttribArray(colorAttr)
	gl.VertexAttribPointer(colorAttr, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
    */

    texCoordAttr  := program.GetAttributeLocation("texCoord")
	gl.EnableVertexAttribArray(texCoordAttr)
	gl.VertexAttribPointer(texCoordAttr, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	proj   := mgl32.Perspective(mgl32.DegToRad(45.0), float32(800.0/600.0), 0.01, 100.0)
	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	model  := mgl32.Ident4()

    program.Matrix4f("projection", &proj[0])
    program.Matrix4f("camera", &camera[0])
    program.Matrix4f("model", &model[0])

    program.Int32("tex", 0)

    angle := float32(0.0)

    //gl.ClearColor(1.0,1.0,1.0,1.0)
    wnd.SetRenderCallback(func(wnd *window.Window, dt float32) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        program.Use()

		angle += dt
		model = mgl32.HomogRotate3D(angle, mgl32.Vec3{0, 1, 0})
        program.Matrix4f("model", &model[0])

        cube.Bind()

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, tx.Id)

        cube.Draw()
    })

    wnd.Loop()
}
