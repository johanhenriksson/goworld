package main

import (
    "github.com/johanhenriksson/goworld/game"
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"

    opensimplex "github.com/ojrac/opensimplex-go"
)

const (
    WIDTH = 1280
    HEIGHT = 800
    WIREFRAME = true
)

func main() {
    app := engine.NewApplication("voxels", WIDTH, HEIGHT)

    /* create a camera */
    app.Scene.Camera = engine.CreateCamera(-3,10,-3, WIDTH, HEIGHT, 65.0, 0.1, 500.0)
    app.Scene.Camera.Transform.Rotation[1] = 130.0

    /* test voxel chunk */
    tilesetMat := render.LoadMaterial(app.Render.Geometry.Shader, "assets/materials/tileset.json")
    tileset := game.CreateTileset(tilesetMat)

    chk := generateChunk(16, tileset)
    chk.Compute()

    obj := engine.NewObject(0,0,0)
    obj.Attach(chk)

    /* attach to scene */
    app.Scene.Add(obj)

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

    bufferWindow("Diffuse", app.Render.Geometry.Buffer.Diffuse, 30, 30)
    bufferWindow("Normal", app.Render.Geometry.Buffer.Normal, 30, 340)

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