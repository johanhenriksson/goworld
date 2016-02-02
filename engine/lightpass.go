package engine

import (
    "github.com/go-gl/gl/v4.1-core/gl"
    //mgl "github.com/go-gl/mathgl/mgl32"
    "github.com/johanhenriksson/goworld/render"
)

type LightPass struct {
    Material *render.Material
    quad     *render.RenderQuad
}

func NewLightPass(input *render.GeometryBuffer) *LightPass {
    /* use a virtual material to help with vertex attributes and textures */
    mat := render.CreateMaterial(render.CompileVFShader("/assets/shaders/voxel_light_pass"))

    /* we're going to render a simple quad, so we input
     * position and texture coordinates */
    mat.AddDescriptor("position", gl.FLOAT, 3, 20, 0, false)
    mat.AddDescriptor("texcoord", gl.FLOAT, 2, 20, 12, false)

    /* the shader uses 3 textures from the geometry frame buffer.
     * they are previously rendered in the geometry pass. */
    mat.AddTexture("tex_diffuse", input.Diffuse)
    mat.AddTexture("tex_normal",  input.Normal)
    mat.AddTexture("tex_depth",   input.Depth)

    /* create a render quad */
    quad := render.NewRenderQuad()
    /* set up vertex attribute pointers */
    mat.SetupVertexPointers()

    p := &LightPass {
        Material: mat,
        quad: quad,
    }
    return p
}

func (p *LightPass) DrawPass(scene *Scene) {
    /* disable depth masking so that multiple lights can be drawn */
    gl.DepthMask(false)

    /* use light pass shader */
    p.Material.Use()
    shader := p.Material.Shader

    /* compute camera view projection inverse */
    vp := scene.Camera.Projection.Mul4(scene.Camera.View)
    vp_inv := vp.Inv()
    shader.Matrix4f("cameraInverse", &vp_inv[0])

    /* clear */
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    /* set blending mode to additive */
    gl.BlendFunc(gl.ONE, gl.ONE)

    /* draw lights */
    lights := scene.FindLights()
    for _, light := range lights {
        /* shadow pass */

        /* set light uniform attributes */
        shader.Vec3("light.Position", &light.Position)
        shader.Vec3("light.Color",    &light.Color)
        shader.Float("light.Range",    light.Range)
        shader.Float("light.attenuation.Constant",  light.Attenuation.Constant)
        shader.Float("light.attenuation.Linear",    light.Attenuation.Linear)
        shader.Float("light.attenuation.Quadratic", light.Attenuation.Quadratic)

        /* render light */
        p.quad.Draw()
    }

    /* reset GL state */
    gl.DepthMask(true)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}