package engine

import (
    "github.com/go-gl/gl/v4.1-core/gl"
    //mgl "github.com/go-gl/mathgl/mgl32"
    "github.com/johanhenriksson/goworld/render"
)

type Renderer struct {
    Width       int32
    Height      int32
    Geometry    *GeometryPass
    Lights      *LightPass
    Scene       *Scene
    time        float32
}

func NewRenderer(width, height int32, scene *Scene) *Renderer {
    gpass := NewGeometryPass(width, height, render.CompileVFShader("/assets/shaders/voxel_geom_pass"))
    lpass := NewLightPass(gpass.Buffer, render.CompileVFShader("/assets/shaders/voxel_light_pass"))

    r := &Renderer {
        Width: width,
        Height: height,
        Geometry: gpass,
        Lights: lpass,
        Scene: scene,
    }

    /* Enable blending */
    gl.Enable(gl.BLEND);
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
    gl.ClearColor(0.0, 0.0, 0.0, 1.0)

    return r
}

func (r *Renderer) Draw() {
    /* Geometry Pass */
    r.Geometry.Draw(r.Scene)

    /* Lighting Pass */
    r.Lights.Draw(r.Scene)
}

func (r *Renderer) Update(dt float32) {
    r.time += dt
    r.Scene.Update(dt)
}
