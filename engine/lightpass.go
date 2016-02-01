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
    mat.Setup()

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

    /* draw lights */
    lights := scene.FindLights()
    for i, light := range lights {
        /* change blending mode on the second element.
           we want the first one to clear */
        if i == 1 {
            // set blending mode to 1 + 1
            gl.BlendFunc(gl.ONE, gl.ONE)
        }

        /* set light uniform attributes */
        p.Shader.Vec3("light.Position", &light.Position)
        p.Shader.Vec3("light.Color", &light.Color)
        p.Shader.Float("light.Range", light.Range)

        p.quad.Draw()
    }

    /* reset GL state */
    gl.DepthMask(true)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}
