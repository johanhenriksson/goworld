package main

import (
    "fmt"
    "github.com/johanhenriksson/goworld/game"
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/physics"

    //mgl "github.com/go-gl/mathgl/mgl32"
    opensimplex "github.com/ojrac/opensimplex-go"
)

const (
    WIDTH = 1280
    HEIGHT = 800
)

func main() {
    app := engine.NewApplication("voxels", WIDTH, HEIGHT)

    /* Setup deferred rendering */
    geom_pass := engine.NewGeometryPass(WIDTH, HEIGHT)
    light_pass := engine.NewLightPass(geom_pass.Buffer)
    app.Render.Append("geometry", geom_pass)
    app.Render.Append("light", light_pass)
    // UI as a render pass?

    /* create a camera */
    app.Scene.Camera = engine.CreateCamera(-3,2,-3, WIDTH, HEIGHT, 65.0, 0.1, 500.0)
    app.Scene.Camera.Transform.Rotation[1] = 130.0

    /* test voxel chunk */
    tileset := game.CreateTileset()
    chk := generateChunk(1, tileset)
    chk.Compute()
    geom_pass.Material.SetupVertexPointers()

    obj := engine.NewObject(0,0,0)
    obj.Attach(chk)
    app.Scene.Add(obj)

    obj2 := engine.NewObject(0.35,3,0)
    obj2.Attach(chk)
    app.Scene.Add(obj2)

    //cam_ray := space.NewRay(10)

    fmt.Println("goworld")
    w := physics.NewWorld()
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
        obj2.Transform.Position = rb2.Position()
        obj2.Transform.Rotation = rb2.Rotation()

        //cam_ray.SetPosDir(toOdeVec3(app.Scene.Camera.Position), toOdeVec3(app.Scene.Camera.Forward))

        w.Update()
    })

    app.Run()
}

func generateChunk(size int, tileset *game.Tileset) *game.Chunk {
    /* Define voxels */
    grass := &game.Voxel{
        Xp: tileset.GetId(4, 0),
        Xn: tileset.GetId(4, 0),
        Yp: tileset.GetId(3, 0),
        Yn: tileset.GetId(2, 0),
        Zp: tileset.GetId(4, 0),
        Zn: tileset.GetId(4, 0),
    }
    rock := &game.Voxel{
        Xp: tileset.GetId(2, 0),
        Xn: tileset.GetId(2, 0),
        Yp: tileset.GetId(2, 0),
        Yn: tileset.GetId(2, 0),
        Zp: tileset.GetId(2, 0),
        Zn: tileset.GetId(2, 0),
    }

    /* Fill chunk with voxels */
    f := 1.0 / 5
    chk := game.CreateChunk(size, tileset)
    simplex := opensimplex.NewWithSeed(1000)
    for z := 0; z < size; z++ {
        for y := 0; y < size; y++ {
            for x := 0; x < size; x++ {
                fx, fy, fz := float64(x) * f, float64(y) * f, float64(z) * f
                v := simplex.Eval3(fx, fy, fz)
                var vtype *game.Voxel = nil
                if y <= size / 2 {
                    vtype = grass
                }
                if v > 0.0 {
                    vtype = rock
                }
                chk.Set(x, y, z, vtype)
            }
        }
    }

    return chk
}