package engine

import (
    //"github.com/go-gl/gl/v4.1-core/gl"
    mgl "github.com/go-gl/mathgl/mgl32"
    "github.com/johanhenriksson/goworld/render"
)

type Component interface {
    Update(float32)
    Draw(render.DrawArgs)
}

type Renderer struct {
    Width       int32
    Height      int32
    Geometry    *GeometryPass
    Lights      *LightPass
    Scene       *Scene
}

func NewRenderer(width, height int32) *Renderer {
    gpass := NewGeometryPass(width, height, render.CompileVFShader("/assets/shaders/voxel_geom_pass"))
    lpass := NewLightPass(gpass, render.CompileVFShader("/assets/shaders/voxel_light_pass"))

    r := &Renderer {
        Width: width,
        Height: height,
        Geometry: gpass,
        Lights: lpass,
        Scene: NewScene(),
    }
    return r
}

func (r *Renderer) Draw() {
    r.Geometry.Draw(r.Scene)
    r.Lights.Draw(r.Scene)
}

type Light struct {
    Position mgl.Vec3
    Color    mgl.Vec3
    Range float32
}
