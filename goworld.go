package main

import (
    "fmt"
    "math"
    "github.com/johanhenriksson/goworld/game"
    "github.com/johanhenriksson/goworld/engine"
    "github.com/johanhenriksson/goworld/render"

    "github.com/ianremmler/ode"
    mgl "github.com/go-gl/mathgl/mgl32"

    opensimplex "github.com/ojrac/opensimplex-go"
//    "github.com/johanhenriksson/goworld/physics"
    "github.com/johanhenriksson/goworld/physics"
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

    obj2 := engine.NewObject(2,0,0)
    obj2.Attach(chk)
    app.Scene.Add(obj2)

    // physics
    ode.Init(0, ode.AllAFlag)
    side := 1.0
    world := ode.NewWorld()
    space := ode.NilSpace().NewSimpleSpace()

    box1 := world.NewBody()
    box1.SetPosition(ode.V3(0, 1.5, 0))
    mass := ode.NewMass()
    mass.SetBox(1, ode.V3(side, side, side))
    mass.Adjust(1)
    box1.SetMass(mass)
    box1_col := space.NewBox(ode.V3(side, side, side))
    box1_col.SetBody(box1)

    box2 := world.NewBody()
    box2.SetPosition(ode.V3(0, 0.2, 0))
    mass2 := ode.NewMass()
    mass2.SetBox(1, ode.V3(side, side, side))
    mass2.Adjust(1)
    box2.SetMass(mass2)
    box2_col := space.NewBox(ode.V3(side, side, side))
    box2_col.SetBody(box2)

    ctGrp := ode.NewJointGroup(1000000)

    world.SetGravity(ode.V3(0,-1,0))
    world.SetCFM(1.0e-5)
    world.SetERP(0.2)
    world.SetContactSurfaceLayer(0.001)
    world.SetContactMaxCorrectingVelocity(0.9)
    world.SetAutoDisable(true)
    space.NewPlane(ode.V4(0,1,0,0))

    //cam_ray := space.NewRay(10)

    fmt.Println("goworld")
    //rb := physics.NewRigidBody(5)
    //fmt.Println(rb)


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
        obj.Transform.Position = physics.FromOdeVec3(box1.Position()).Sub(mgl.Vec3{0.5,0.5,0.5})
        obj.Transform.Rotation = physics.FromOdeRotation(box1.Rotation())
        obj2.Transform.Position = physics.FromOdeVec3(box2.Position()).Sub(mgl.Vec3{0.5,0.5,0.5})
        obj2.Transform.Rotation = physics.FromOdeRotation(box2.Rotation())

        //cam_ray.SetPosDir(toOdeVec3(app.Scene.Camera.Position), toOdeVec3(app.Scene.Camera.Forward))

        space.Collide(0, func(data interface{}, obj1, obj2 ode.Geom) {
            body1, body2 := obj1.Body(), obj2.Body()
            cts := obj1.Collide(obj2, 1, 0)
            for _, ct := range cts {
                contact := ode.NewContact()
                contact.Surface.Mode = ode.BounceCtParam | ode.SoftCFMCtParam;
                contact.Surface.Mu = math.Inf(1)
                contact.Surface.Mu2 = 0;
                contact.Surface.Bounce = 0.01;
                contact.Surface.BounceVel = 0.1;

                contact.Geom = ct
                ct := world.NewContactJoint(ctGrp, contact)
                ct.Attach(body1, body2)
            }
        })
        world.QuickStep(0.01)
        ctGrp.Empty()
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