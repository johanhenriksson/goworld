package engine

import (
    "github.com/johanhenriksson/goworld/assets"
    "github.com/johanhenriksson/goworld/render"
    //"github.com/go-gl/gl/v4.1-core/gl"
)

type LinePass struct {
    Material *render.Material
}

/* Sets up a geometry pass.
 * A geometry buffer of the given bufferWidth x bufferHeight will be created automatically */
func NewLinePass() *LinePass {
    mat := assets.GetMaterialCached("lines")
    p := &LinePass {
        Material: mat,
    }
    return p
}

func (p *LinePass) DrawPass(scene *Scene) {
    //gl.DepthMask(false)

    p.Material.Use()
    shader := p.Material.Shader
    shader.Matrix4f("view", &scene.Camera.View[0])
    shader.Matrix4f("projection", &scene.Camera.Projection[0])

    /* Draw scene */
    scene.Draw("lines", shader)

    //gl.DepthMask(true)
}