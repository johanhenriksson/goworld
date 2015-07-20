package main

import (
    "fmt"
    "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"

    "github.com/johanhenriksson/goworld/window"
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"
)

func main() {
    wnd := window.Create("voxels", 800, 600)

    /* Perspective & camera */
    /*
	proj   := mgl32.Perspective(mgl32.DegToRad(45.0), float32(800.0/600.0), 0.01, 100.0)
	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
    */
	model  := mgl32.Ident4()

    /* Shader setup */
    program := render.CompileVFShader("assets/shaders/3d_voxel")
    program.Use()

    cam := engine.CreateCamera(0,3,5, 800,600, 45.0, 0.1, 100.0)
    program.Matrix4f("projection", &cam.Projection[0])
    program.Matrix4f("camera", &cam.View[0])

    program.Matrix4f("model", &model[0])

    ttx, _ := render.LoadTexture("assets/tileset.png")
    tilesetMat := render.CreateMaterial(program)
    tilesetMat.AddDescriptor("vertex", gl.UNSIGNED_BYTE, 3, 8, 0, false)
    tilesetMat.AddDescriptor("normal", gl.BYTE, 3, 8, 3, false)
    tilesetMat.AddDescriptor("tile", gl.UNSIGNED_BYTE, 2, 8, 6, false)
    tilesetMat.AddTexture(0, ttx)

    tileset := engine.CreateTileset(tilesetMat)
    fmt.Println("Tileset", tileset.Width, "x", tileset.Height)


    /* Voxel mesh */

    /* Define a gress tile */
    grass := &engine.Voxel {
        Xp: tileset.GetId(4, 0),
        Xn: tileset.GetId(4, 0),
        Yp: tileset.GetId(3, 0),
        Yn: tileset.GetId(2, 0),
        Zp: tileset.GetId(4, 0),
        Zn: tileset.GetId(4, 0),
    }

    /* Fill chunk with voxels */
    chk := engine.CreateChunk(32, tileset)
    for i := 0; i < 32*32*8; i++ {
        chk.Data[i] = grass
    }

    /* Compute mesh */
    vmesh := chk.Compute()

    gl.ClearColor(1,1,1,0)

    /* Render loop */
    angle := float32(0.0)
    wnd.SetRenderCallback(func(wnd *window.Window, dt float32) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        program.Use()

		angle += dt
		model = mgl32.HomogRotate3D(angle, mgl32.Vec3{0, 1, 0})
        program.Matrix4f("camera", &cam.View[0])
        program.Matrix4f("model", &model[0])

        vmesh.Render()
    })

    wnd.SetUpdateCallback(func(dt float32) {
        cam.Update(dt)
    })

    wnd.Loop()
}
