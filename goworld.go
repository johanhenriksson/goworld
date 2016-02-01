package main

import (
    "fmt"
    "time"
    "github.com/go-gl/gl/v4.1-core/gl"
    //mgl "github.com/go-gl/mathgl/mgl32"

    "github.com/johanhenriksson/goworld/game"
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/ui"

    opensimplex "github.com/ojrac/opensimplex-go"
)

const (
    WIDTH = 1280
    HEIGHT = 800
    WIREFRAME = true
)

func main() {
    wnd := engine.CreateWindow("voxels", WIDTH, HEIGHT)

    gl.ClearColor(0.0,0.0,0.0,1)

    /* Main camera */
    cam := engine.CreateCamera(5,2,5, WIDTH, HEIGHT, 65.0, 0.1, 500.0)

    /* Scene */
    scene := engine.NewScene(cam)

    /* Renderer */
    rnd := engine.NewRenderer(WIDTH, HEIGHT, scene)

    /* UI Manager */
    uimgr := ui.NewManager(wnd)

    fmt.Println("goworld running")

    tilesetMat := render.LoadMaterial(rnd.Geometry.Shader, "assets/materials/tileset.json")
    tileset := game.CreateTileset(tilesetMat)

    chk := generateChunk(16, tileset)
    obj := engine.NewObject(0,0,0)
    obj.Attach(chk)
    chk.Compute()
    rnd.Scene.Add(obj)

    // buffer display window
    bufferWindow := func(title string, texture *render.Texture, x, y float32) {
        win_color := render.Color{0.15, 0.15, 0.15, 0.8}
        text_color := render.Color{1,1,1,1}

        win   := uimgr.NewRect(win_color, x, y, 250, 280, -10)
        label := uimgr.NewText(title, text_color, 0, 0, -21)
        win.Append(label)
        img := uimgr.NewImage(texture, 0, 30, 250, 250, -20)
        img.Quad.FlipY()
        win.Append(img)
        uimgr.Append(win)
    }

    bufferWindow("Diffuse", rnd.Geometry.Buffer.Diffuse, 30, 30)
    bufferWindow("Normal", rnd.Geometry.Buffer.Normal, 30, 340)

    /* Render loop */
    wnd.SetRenderCallback(func(wnd *engine.Window, dt float32) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        /* render scene */
        rnd.Draw()

        /* draw user interface */
        uimgr.Draw()

        time.Sleep(2 * time.Millisecond)
    })

    wnd.SetUpdateCallback(rnd.Update)

    wnd.Loop()
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