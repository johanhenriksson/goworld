package engine

import (
    "github.com/go-gl/gl/v4.1-core/gl"
    mgl "github.com/go-gl/mathgl/mgl32"
    "github.com/johanhenriksson/goworld/render"
)

type GeometryPass struct {
    Buffer *render.GeometryBuffer
    Shader *render.ShaderProgram
}

func NewGeometryPass(bufferWidth, bufferHeight int32, shader *render.ShaderProgram) *GeometryPass {
    p := &GeometryPass {
        Buffer: render.CreateGeometryBuffer(bufferWidth, bufferHeight),
        Shader: shader,
    }
    return p
}

func (p *GeometryPass) Draw(scene *Scene) {
    p.Buffer.Bind()
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    cam := scene.Camera

    p.Shader.Use()
    m := mgl.Ident4()
    p.Shader.Matrix4f("model", &m[0])
    p.Shader.Matrix4f("camera", &cam.View[0])
    p.Shader.Matrix4f("projection", &cam.Projection[0])

    scene.Draw(p.Shader)

    p.Buffer.Unbind()
}
