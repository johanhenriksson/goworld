package main

import (
    "fmt"
    "math"
    "github.com/go-gl/gl/v4.1-core/gl"
    mgl "github.com/go-gl/mathgl/mgl32"

    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/geometry"
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/ui"

    opensimplex "github.com/ojrac/opensimplex-go"
)

const (
    WIDTH = 1280
    HEIGHT = 800
)

func main() {
    wnd := engine.CreateWindow("voxels", WIDTH, HEIGHT)
    cam := engine.CreateCamera(5,2,5, WIDTH, HEIGHT, 65.0, 0.1, 1000.0)
    //gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

    /* Enable blending */
    gl.Enable(gl.BLEND);
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);

    uimgr := ui.NewManager(wnd)
    rect := uimgr.NewRect(ui.NewColor(0.5, 1, 0.7, 0.5), 120, 80, 800, 500, -10)
    kitten, err := render.LoadTexture("/assets/kitten.png")
    if err != nil {
        panic(err)
    }
    img := uimgr.NewImage(kitten, 10, 20, 400, 300, -20)
    rect.Append(img)
    uimgr.Append(rect)

    /* Line material */
    lineMat := render.LoadMaterial("assets/materials/lines.json")
    lineProgram := lineMat.Shader

    tilesetMat := render.LoadMaterial("assets/materials/tileset.json")
    program := tilesetMat.Shader
    program.Matrix4f("projection", &cam.Projection[0])

    tileset := engine.CreateTileset(tilesetMat)


    /* Define voxels */
    grass := &engine.Voxel {
        Xp: tileset.GetId(4, 0),
        Xn: tileset.GetId(4, 0),
        Yp: tileset.GetId(3, 0),
        Yn: tileset.GetId(2, 0),
        Zp: tileset.GetId(4, 0),
        Zn: tileset.GetId(4, 0),
    }
    rock := &engine.Voxel {
        Xp: tileset.GetId(2, 0),
        Xn: tileset.GetId(2, 0),
        Yp: tileset.GetId(2, 0),
        Yn: tileset.GetId(2, 0),
        Zp: tileset.GetId(2, 0),
        Zn: tileset.GetId(2, 0),
    }


    /* Fill chunk with voxels */
    size := 16
    f := 1.0 / 5
    chk := engine.CreateChunk(size, tileset)
    simplex := opensimplex.NewWithSeed(1000)
    for z := 0; z < size; z++ {
        for y := 0; y < size; y++ {
            for x := 0; x < size; x++ {
                fx, fy, fz := float64(x) * f, float64(y) * f, float64(z) * f
                v := simplex.Eval3(fx, fy, fz)
                var vtype *engine.Voxel = nil
                if y <= size/2 {
                    vtype = grass
                }
                if v > 0.0 {
                    vtype = rock
                }
                chk.Set(x,y,z, vtype)
            }
        }
    }

    transf := engine.CreateTransform(0,0,0)

    /* Lines */
    lines := geometry.CreateLines(lineMat)
    lines.Box(0,0,0,16,16,16,0.5,1,0.5,1)
    lines.Compute()
    lineProgram.Use()
    lineProgram.Matrix4f("projection", &cam.Projection[0])
    lineProgram.Matrix4f("model", &transf.Matrix[0])

    /* Compute mesh */
    vmesh := chk.Compute()
    program.Use()
    program.Matrix4f("model", &transf.Matrix[0])
    program.Vec3("lightPos", &mgl.Vec3{ 8,15,8 })
    program.Float("lightIntensity", 250.0)
    program.Float("ambient", 0.6)

    gl.ClearColor(0,0,0,1)

    /* Render loop */
    wnd.SetRenderCallback(func(wnd *engine.Window, dt float32) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

        program.Use()
        program.Matrix4f("camera", &cam.View[0])
        program.Vec3("cameraPos", &cam.Transform.Position)

        vmesh.Render()

        lineProgram.Use()
        lineProgram.Matrix4f("view", &cam.View[0])
        lines.Render()

        uimgr.Draw()

    })

    shoot := false
    wnd.SetUpdateCallback(func(dt float32) {
        if engine.KeyDown(engine.KeyF) {
            if !shoot {
                pos := cam.Unproject(1280 / 2, 800 / 2)
                x,y,z := VoxelCoord(cam.Forward, pos)
                fmt.Println(x,y,z)
                chk.Set(x,y,z,grass)
                vmesh = chk.Compute()
                shoot = true
            }
        } else {
            shoot = false
        }
        cam.Update(dt)
    })

    wnd.Loop()
}

func dti(val float32) float32 {
  return float32(math.Abs(float64(val - Round(val))));
}

func Round(f float32) float32 {
    return float32(math.Floor(float64(f + .5)))
}

func VoxelCoord(forward mgl.Vec3, coord mgl.Vec3) (int, int, int) {
    nx := int(coord[0]);
    ny := int(coord[1]);
    nz := int(coord[2]);

    /* find the coordinate that is closer to an integer value */
    /* x < y? */
    if dti(coord[0]) < dti(coord[1]) {
        /* x is less than y */
        /* x < z? */
        if dti(coord[0]) < dti(coord[2]) {
            /* x is closer */
            if forward[0] > 0 {
                /* we are looking to the right */
                fmt.Println("X closest, looking along X+")
                //nx--;
            } else {
                nx++;
                fmt.Println("X closest, looking along X-")
            }
        } else {
            /* z is closer */
            if forward[2] > 0 {
                /* we are looking along z+ */
                fmt.Println("1 Z closest, looking along Z+")
                //nz--
            } else {
                //nz++
                fmt.Println("1 Z closest, looking along Z-")
            }
        }
    } else {
        /* x > y */
        /* is y closer than z? */
        if dti(coord[1]) < dti(coord[2]) {
            /* y is closer! */
            if forward[1] > 0 {
                /* we are looking up */
                fmt.Println("Y closest, looking up")
                ny--
            }
        } else {
            /* z is closer! */
            if forward[2] > 0 {
                /* looking along z+ */
                fmt.Println("1 Z closest, looking along Z+")
                //nz--
            } else {
                //nz++
                fmt.Println("1 Z closest, looking along Z-")
            }
        }
    }
    return nx, ny, nz
}
