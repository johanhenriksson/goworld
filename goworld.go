package main

import (
    "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

    "github.com/johanhenriksson/goworld/window"
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/geometry"
)

func main() {
    wnd := window.Create("Hello World", 800, 600)

    /* Perspective & camera */
	proj   := mgl32.Perspective(mgl32.DegToRad(45.0), float32(800.0/600.0), 0.01, 100.0)
	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	model  := mgl32.Ident4()

    /* Shader setup */
    program := render.CreateProgram()
    program.Attach(render.VertexShader("assets/shaders/3d_textured.vs.glsl"))
    program.Attach(render.FragmentShader("assets/shaders/3d_textured.fs.glsl"))
    program.Link()
    program.Use()

    program.Matrix4f("projection", &proj[0])
    program.Matrix4f("camera", &camera[0])
    program.Matrix4f("model", &model[0])

    /* Mesh */
    tx, _ := render.LoadTexture("assets/square.png")

    mat := render.CreateMaterial(program, "vertex", "normal", "texCoord")
    mat.AddTexture(gl.TEXTURE0, tx)

    /* Set up cube */
    cube := geometry.Cube()
    mesh := engine.CreateMesh(cube, mat)

    /* Render loop */
    angle := float32(0.0)
    wnd.SetRenderCallback(func(wnd *window.Window, dt float32) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        program.Use()

		angle += dt
		model = mgl32.HomogRotate3D(angle, mgl32.Vec3{0, 1, 0})
        program.Matrix4f("model", &model[0])

        mesh.Render()
    })

    wnd.Loop()
}
