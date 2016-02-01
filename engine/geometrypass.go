package engine

import (
    //"github.com/go-gl/gl/v4.1-core/gl"
    //mgl "github.com/go-gl/mathgl/mgl32"
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
    p.Buffer.Clear()

    p.Shader.Use()
    p.Shader.Matrix4f("camera", &scene.Camera.View[0])
    p.Shader.Matrix4f("projection", &scene.Camera.Projection[0])

    scene.Draw(p.Shader)

    p.Buffer.Unbind()
}
