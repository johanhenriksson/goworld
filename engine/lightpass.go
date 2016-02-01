package engine

import (
    "github.com/go-gl/gl/v4.1-core/gl"
    //mgl "github.com/go-gl/mathgl/mgl32"
    "github.com/johanhenriksson/goworld/render"
)

type LightPass struct {
    Input   *render.GeometryBuffer
    Shader  *render.ShaderProgram
    Lights  []Light

    mat     *render.Material
    quad    *render.RenderQuad
}

func NewLightPass(geomPass *GeometryPass, shader *render.ShaderProgram) *LightPass {
    /* set up some shader defaults */
    shader.Use()
    shader.Float("light.attenuation.Constant", 0.02);
    shader.Float("light.attenuation.Linear", 0.0);
    shader.Float("light.attenuation.Quadratic", 1.0);

    mat := render.CreateMaterial(shader)
    /* we're going to render a simple quad, so we input
     * position and texture coordinates */
    mat.AddDescriptor("position", gl.FLOAT, 3, 20, 0, false)
    mat.AddDescriptor("texcoord", gl.FLOAT, 2, 20, 12, false)

    /* the shader uses 3 textures - the geometry frame buffer
     * textures previously rendered in the geometry pass. */
    mat.AddTexture("tex_diffuse", geomPass.Buffer.Diffuse)
    mat.AddTexture("tex_normal", geomPass.Buffer.Normal)
    mat.AddTexture("tex_depth", geomPass.Buffer.Depth)

    /* create a render quad */
    quad := render.NewRenderQuad()
    /* set up vertex attribute pointers */
    mat.SetupVertexPointers()

    p := &LightPass {
        Input: geomPass.Buffer,
        Shader: shader,
        quad: quad,
        mat: mat,
    }
    return p
}

func (p *LightPass) Draw(scene *Scene) {
    /* disable depth masking so that multiple lights can be drawn */
    gl.DepthMask(false)

    /* use light pass shader */
    p.mat.Use()

    /* compute camera view projection inverse */
    vp := scene.Camera.Projection.Mul4(scene.Camera.View)
    vp_inv := vp.Inv()
    p.Shader.Matrix4f("cameraInverse", &vp_inv[0])

    /* clear */
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    /* set blending mode to additive */
    gl.BlendFunc(gl.ONE, gl.ONE)

    /* draw lights */
    lights := scene.FindLights()
    for _, light := range lights {
        /* set light uniform attributes */
        p.Shader.Vec3("light.Position", &light.Position)
        p.Shader.Vec3("light.Color", &light.Color)
        p.Shader.Float("light.Range", light.Range)

        /* render light */
        p.quad.Draw()
    }

    /* reset GL state */
    gl.DepthMask(true)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
