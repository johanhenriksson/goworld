package main

import (
    "fmt"
    "github.com/johanhenriksson/goworld/game"
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"
    //"github.com/johanhenriksson/goworld/geometry"

    //mgl "github.com/go-gl/mathgl/mgl32"
    opensimplex "github.com/ojrac/opensimplex-go"
    //"github.com/go-gl/mathgl/mgl32"
)

const (
    WIDTH = 1280
    HEIGHT = 800
)

func main() {
    app := engine.NewApplication("voxels", WIDTH, HEIGHT)

    /* Setup deferred rendering */
    geom_pass := engine.NewGeometryPass(WIDTH, HEIGHT)
    app.Render.Append("geometry", geom_pass)
    app.Render.Append("light", engine.NewLightPass(geom_pass.Buffer))
    app.Render.Append("lines", engine.NewLinePass())
    // UI as a render pass?

    /* create a camera */
    app.Scene.Camera = engine.CreateCamera(-3,2,-3, WIDTH, HEIGHT, 65.0, 0.1, 500.0)
    app.Scene.Camera.Transform.Rotation[1] = 130.0

    obj := app.Scene.NewObject(3,0,3)

    /* test voxel chunk */
    tileset := game.CreateTileset()
    chk := game.NewChunk(obj, 20, tileset)
    chk.Compute()

    //game.NewPlacementGrid(obj)

    //app.Scene.Add(obj)

    obj2 := app.Scene.NewObject(2,1,0)
    chk2 := game.NewColorChunk(obj2, 16)
    generateChunk(chk2) // populate with random data
    chk2.Set(0,0,0, &game.ColorVoxel{ R:255, G:0, B:0 })
    chk2.Set(1,0,0, &game.ColorVoxel{ R:0, G:255, B:0 })
    chk2.Set(2,0,0, &game.ColorVoxel{ R:0, G:0, B:255 })
    chk2.Compute()
    geom_pass.Material.SetupVertexPointers()
    app.Scene.Add(obj2)

    game.NewPlacementGrid(obj2)

    //cam_ray := space.NewRay(10)

    fmt.Println("goworld")
    w := app.Scene.World
    w.NewPlane(0,1,0,0)
    rb1 := w.NewRigidBox(5, 1, 1, 1)
    rb1.SetPosition(obj.Transform.Position)
    rb2 := w.NewRigidBox(5, 1, 1, 1)
    rb2.SetPosition(obj2.Transform.Position)
    fmt.Println(rb1)




    // buffer display window
    bufferWindow := func(title string, texture *render.Texture, x, y float32) {
        win_color := render.Color{0.15, 0.15, 0.15, 0.8}
        text_color := render.Color{1, 1, 1, 1}

        win   := app.UI.NewRect(win_color, x, y, 250, 280, -10)
        label := app.UI.NewText(title, text_color, 0, 0, -21)
        img   := app.UI.NewImage(texture, 0, 30, 250, 250, -20)
        img.Quad.FlipY()

        win.Append(img)
        win.Append(label)

        /* attach UI element */
        app.UI.Append(win)
    }

    bufferWindow("Diffuse", geom_pass.Buffer.Diffuse, 30, 30)
    bufferWindow("Normal", geom_pass.Buffer.Normal, 30, 340)

    /* Render loop */
    app.Window.SetRenderCallback(func(wnd *engine.Window, dt float32) {
        /* render scene */
        app.Render.Draw()

        /* draw user interface */
        app.UI.Draw()

        // update position
        //fmt.Println(box1.Position())
        //fmt.Println(app.Scene.Camera.Forward)
        obj.Transform.Position = rb1.Position()
        obj.Transform.Rotation = rb1.Rotation()
        /*
        obj2.Transform.Position = rb2.Position()
        obj2.Transform.Rotation = rb2.Rotation()
        */

        //cam_ray.SetPosDir(toOdeVec3(app.Scene.Camera.Position), toOdeVec3(app.Scene.Camera.Forward))

        if engine.KeyReleased(engine.KeyF) {
            fmt.Println("raycast")
            w.Raycast(10, app.Scene.Camera.Position, app.Scene.Camera.Forward)
        }
    })

    app.Run()
}

func generateChunk(chk *game.ColorChunk) {
    /* Define voxels */
    rock2 := &game.ColorVoxel{
        R: 200,
        G: 179,
        B: 112,
    }
    rock := &game.ColorVoxel{
        R: 141,
        G: 119,
        B: 72,
    }
    grass := &game.ColorVoxel{
        R: 88,
        G: 132,
        B: 69,
    }

    /* Fill chunk with voxels */
    f := 1.0 / 5
    size := chk.Size
    simplex := opensimplex.NewWithSeed(1000)
    for z := 0; z < size; z++ {
        for y := 0; y < size; y++ {
            for x := 0; x < size; x++ {
                fx, fy, fz := float64(x) * f, float64(y) * f, float64(z) * f
                v := simplex.Eval3(fx, fy, fz)
                var vtype *game.ColorVoxel = nil
                if y < size / 2 {
                    vtype = rock2
                }
                if y == size / 2 {
                    vtype = grass
                }
                if v > 0.0 {
                    vtype = rock
                }
                chk.Set(x, y, z, vtype)
            }
        }
    }
}
