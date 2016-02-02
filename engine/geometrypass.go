package engine

import (
    "github.com/johanhenriksson/goworld/render"
)

type GeometryPass struct {
    Buffer *render.GeometryBuffer
    Material *render.Material
}

/* Sets up a geometry pass.
 * A geometry buffer of the given bufferWidth x bufferHeight will be created automatically */
func NewGeometryPass(bufferWidth, bufferHeight int32) *GeometryPass {
    shader := render.CompileVFShader("assets/shaders/voxel_geom_pass")
    mat := render.LoadMaterial(shader, "assets/materials/tileset")
    p := &GeometryPass {
        Buffer: render.CreateGeometryBuffer(bufferWidth, bufferHeight),
        Material: mat,
    }
    return p
}

func (p *GeometryPass) DrawPass(scene *Scene) {
    p.Buffer.Bind()
    p.Buffer.Clear()

    p.Material.Use()
    shader := p.Material.Shader
    shader.Matrix4f("camera", &scene.Camera.View[0])
    shader.Matrix4f("projection", &scene.Camera.Projection[0])

    /* Draw scene */
    scene.Draw(shader)

    p.Buffer.Unbind()
}