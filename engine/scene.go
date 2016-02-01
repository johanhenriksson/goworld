package engine

import (
    "github.com/johanhenriksson/goworld/render"
    mgl "github.com/go-gl/mathgl/mgl32"
)

/* Scene Graph */
type Scene struct {
    /* Active camera */
    Camera      *Camera

    /* Root Objects */
    Objects     []*Object

    /* temporary: list of all lights in the scene */
    lights      []Light
}

func NewScene() *Scene {
    s := &Scene {
        Camera: nil,
        Objects: []*Object { },
        lights: []Light {
            /* temporary: test light */
            Light {
                Attenuation: Attenuation {
                    Constant: 0.01,
                    Linear: 0,
                    Quadratic: 1.0,
                },
                Color: mgl.Vec3 { 1,1,1 },
                Range: 1,
            },
        },
    }

    /* add a few more test lights */
    for i := 0; i < 5; i++ {
        s.lights = append(s.lights, Light {
            Attenuation: Attenuation {
                Constant: 0.01,
                Linear: 0,
                Quadratic: 1.0,
            },
            Position: mgl.Vec3 { float32(2*i), 17, 0 },
            Color: mgl.Vec3 { 0.3, 0.3, 1.0 },
            Range: 0.1,
        })
    }
    return s
}

func (s *Scene) Add(object *Object) {
    /* TODO look for lights - maybe not here? */
    s.Objects = append(s.Objects, object)
}

func (s *Scene) Draw(shader *render.ShaderProgram) {
    if s.Camera == nil {
        return
    }

    /* DrawArgs will be copied down recursively into the scene graph.
     * Each object adds its transformation matrix before passing
     * it on to their children */
    args := render.DrawArgs {
        Viewport: s.Camera.View,
        Transform: mgl.Ident4(),
        Shader: shader,
    }

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
    s.lights[0].Position = s.Camera.Position
}

func (s *Scene) FindLights() []Light {
    // todo
    return s.lights
}
