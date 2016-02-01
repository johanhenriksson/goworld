package engine

import (
    "github.com/johanhenriksson/goworld/render"
    mgl "github.com/go-gl/mathgl/mgl32"
)

type Scene struct {
    Camera      *Camera
    Objects     []*Object
    lights      []Light
}

func NewScene(camera *Camera) *Scene {
    s := &Scene {
        Camera: camera,
        Objects: []*Object { },
        lights: []Light {
            Light {
                Color: mgl.Vec3 { 1,1,1 },
                Range: 1,
            },
        },
    }
    for i := 0; i < 5; i++ {
        s.lights = append(s.lights, Light {
            Position: mgl.Vec3 { float32(2*i), 17, 0 },
            Color: mgl.Vec3 { 0.3, 0.3, 1.0 },
            Range: 0.1,
        })
    }
    return s
}

func (s *Scene) Add(object *Object) {
    /* TODO look for lights */
    s.Objects = append(s.Objects, object)
}

func (s *Scene) Draw(shader *render.ShaderProgram) {
    if s.Camera == nil {
        return
    }
    s.lights[0].Position = s.Camera.Position
    args := render.DrawArgs {
        Viewport: s.Camera.View,
        Transform: mgl.Ident4(),
        Shader: shader,
    }
    for _, obj := range s.Objects {
        obj.Draw(args)
    }
}

func (s *Scene) Update(dt float32) {
    for _, obj := range s.Objects {
        obj.Update(dt)
    }
    s.Camera.Update(dt)
}

func (s *Scene) FindLights() []Light {
    return s.lights
}
