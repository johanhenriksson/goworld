package engine

import (
    "github.com/go-gl/gl/v4.1-core/gl"
    //mgl "github.com/go-gl/mathgl/mgl32"
)

type Renderer struct {
    Passes      []RenderPass
    Scene       *Scene
    Width       int32
    Height      int32
}

func NewRenderer(width, height int32, scene *Scene) *Renderer {

    gpass := NewGeometryPass(width, height)
    lpass := NewLightPass(gpass.Buffer)



    r := &Renderer {
        Width: width,
        Height: height,
        Scene: scene,
        Passes: []RenderPass {
            gpass,
            lpass,
        },
    }

    /* Enable blending */
    gl.Enable(gl.BLEND);
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
    gl.ClearColor(0.0, 0.0, 0.0, 1.0)

    return r
}

func (r *Renderer) Draw() {
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    for _, pass := range r.Passes {
        pass.DrawPass(r.Scene)
    }
}

func (r *Renderer) Update(dt float32) {
    r.Scene.Update(dt)
}