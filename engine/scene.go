package engine

import (
    "github.com/johanhenriksson/goworld/render"
    "github.com/johanhenriksson/goworld/physics"

    mgl "github.com/go-gl/mathgl/mgl32"
)


/* Scene Graph */
type Scene struct {
    /* Active camera */
    Camera      *Camera

    /* Root Objects */
    Objects     []*Object

    World       *physics.World

    /* temporary: list of all lights in the scene */
    lights      []Light
}

func NewScene() *Scene {
    s := &Scene {
        Camera: nil,
        Objects: []*Object { },
        World: physics.NewWorld(),

        lights: []Light {
            /* temporary: test light */
            Light {
                Attenuation: Attenuation {
                    Constant: 0.01,
                    Linear: 0,
                    Quadratic: 1.0,
                },
                Color: mgl.Vec3 { 1,1,1 },
                Range: 4,
                Type: PointLight,
            },
        },
    }

    /* add a few more test lights */
    /*
    for i := 0; i < 1; i++ {
        s.lights = append(s.lights, Light {
            Attenuation: Attenuation {
                Constant: 0.01,
                Linear: 0.5,
                Quadratic: 1.5,
            },
            Position: mgl.Vec3 { 0, float32(4*i), 0 },
            Color: mgl.Vec3 { 1.0, 1.0, 0.65 },
            Range: 4,
            Type: PointLight,
        })
    }
    */

    s.lights[0].Position = mgl.Vec3 { -11, 11, -11 }
    s.lights[0].Color = mgl.Vec3 { 0.95, 0.95, 0.96 }
    s.lights[0].Type = DirectionalLight
    s.lights[0].Projection = mgl.Ortho(-32,32,0,64,-32,64)
    //s.lights[2].Position = mgl.Vec3 { 10, 40, 10 }
    //s.lights[2].Range = 10
    //s.lights[0].Type = 2
    return s
}

func (s *Scene) Add(object *Object) {
    /* TODO look for lights - maybe not here? */
    s.Objects = append(s.Objects, object)
}

func (s *Scene) Draw(pass string, shader *render.ShaderProgram) {
    if s.Camera == nil {
        return
    }

    p := s.Camera.Projection
    v := s.Camera.View
    m := mgl.Ident4()
    vp := p.Mul4(v)
    // mvp := vp * m

    /* DrawArgs will be copied down recursively into the scene graph.
     * Each object adds its transformation matrix before passing
     * it on to their children */
    args := render.DrawArgs {
        Projection: p,
        View: v,
        VP: vp,
        MVP: vp,
        Transform: m,

        Pass: pass,
        Shader: shader,
    }

    s.DrawCall(args)
}

func (s *Scene) DrawCall(args render.DrawArgs) {
    /* draw root objects */
    for _, obj := range s.Objects {
        obj.Draw(args)
    }
}

func (s *Scene) Update(dt float32) {
    if s.Camera != nil {
        /* update camera first */
        s.Camera.Update(dt)
    }

    /* update root objects */
    for _, obj := range s.Objects {
        obj.Update(dt)
    }

    /* test: position first light on camera */
    //s.lights[0].Position = s.Camera.Position

    /* physics step */
    s.World.Update()
}

func (s *Scene) FindLights() []Light {
    // todo
    return s.lights
}
