package main

import (
    "github.com/go-gl/gl/v4.1-core/gl"

    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"
)

func main() {
    wnd := engine.CreateWindow("voxels", 1280, 800)

    cam := engine.CreateCamera(5,2,5, 1280,800, 65.0, 0.1, 100.0)

    /* Shader setup */
    program := render.CompileVFShader("assets/shaders/3d_voxel")
    program.Use()
    program.Matrix4f("projection", &cam.Projection[0])
    program.Matrix4f("camera", &cam.View[0])

    /* Tileset Material */
    ttx, _ := render.LoadTexture("assets/tileset.png")
    tilesetMat := render.CreateMaterial(program)
    tilesetMat.AddDescriptor("vertex", gl.UNSIGNED_BYTE, 3, 8, 0, false)
    tilesetMat.AddDescriptor("normal", gl.BYTE, 3, 8, 3, false)
    tilesetMat.AddDescriptor("tile", gl.UNSIGNED_BYTE, 2, 8, 6, false)
    tilesetMat.AddTexture(0, ttx)

    tileset := engine.CreateTileset(tilesetMat)

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
    for i := 0; i < 2*8; i++ {
        chk.Data[i] = grass
    }
    chk.Set(4,1,0, grass)

    /* Compute mesh */
    vmesh := chk.Compute()
    transf := engine.CreateTransform(0,0,0)
    program.Matrix4f("model", &transf.Matrix[0])

    gl.ClearColor(1,1,1,0)

    /* Render loop */
    wnd.SetRenderCallback(func(wnd *engine.Window, dt float32) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        program.Matrix4f("camera", &cam.View[0])

        vmesh.Render()
    })

    wnd.SetUpdateCallback(cam.Update)

    wnd.Loop()
}
