package engine

import (
    "github.com/johanhenriksson/goworld/render"
)

type GeometryPass struct {
    Buffer *render.GeometryBuffer
    Shader *render.ShaderProgram
}

/* Sets up a geometry pass.
 * A geometry buffer of the given bufferWidth x bufferHeight will be created automatically */
func NewGeometryPass(bufferWidth, bufferHeight int32, shader *render.ShaderProgram) *GeometryPass {
    p := &GeometryPass {
        Buffer: render.CreateGeometryBuffer(bufferWidth, bufferHeight),
        Shader: shader,
    }
    return p
}

func (p *GeometryPass) Draw(scene *Scene) {
    p.Buffer.Bind()
    p.Buffer.Clear()

    p.Shader.Use()
    p.Shader.Matrix4f("camera", &scene.Camera.View[0])
    p.Shader.Matrix4f("projection", &scene.Camera.Projection[0])

    /* Draw scene */
    scene.Draw(p.Shader)

    p.Buffer.Unbind()
}