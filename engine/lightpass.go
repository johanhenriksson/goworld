package engine

import (
    "github.com/go-gl/gl/v4.1-core/gl"
    mgl "github.com/go-gl/mathgl/mgl32"
    "github.com/johanhenriksson/goworld/render"
)

type LightPass struct {
    Material *render.Material
    quad     *render.RenderQuad
    Shadows  *ShadowPass
}

func NewLightPass(input *render.GeometryBuffer) *LightPass {
    /* use a virtual material to help with vertex attributes and textures */
    mat := render.CreateMaterial(render.CompileVFShader("/assets/shaders/voxel_light_pass"))

    /* we're going to render a simple quad, so we input
     * position and texture coordinates */
    mat.AddDescriptor("position", gl.FLOAT, 3, 20, 0, false, false)
    mat.AddDescriptor("texcoord", gl.FLOAT, 2, 20, 12, false, false)

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
        Shadows: NewShadowPass(input),
    }
    return p
}

func (p *LightPass) DrawPass(scene *Scene) {
    /* use light pass shader */
    p.Material.Use()
    shader := p.Material.Shader

    /* compute camera view projection inverse */
    vp := scene.Camera.Projection.Mul4(scene.Camera.View)
    vp_inv := vp.Inv()
    shader.Matrix4f("cameraInverse", &vp_inv[0])

    /* clear */
    gl.ClearColor(0.9,0.9,0.9,1.0)
    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

    /* set blending mode to additive */

    gl.DepthMask(false)

    /* draw lights */
    lights := scene.FindLights()
    last := len(lights) - 1

    for i, light := range lights {
        if i == 1 {
            /* first light pass we want the shader to restore the depth buffer
             * then, disable depth masking so that multiple lights can be drawn */
            gl.BlendFunc(gl.ONE, gl.ONE)
        }
        if i == last {
            gl.DepthMask(true)
        }

        /* shadow pass */
        p.Shadows.DrawPass(scene, &light)

        /* use light shader */
        p.Material.SetTexture("tex_shadow", p.Shadows.Output)
        p.Material.Use()

        lp := mgl.Ortho(-150,150, -150,150, -150,150)
        lv := mgl.LookAtV(mgl.Vec3{0,0,0}, light.Position.Normalize(), mgl.Vec3{0,1,0}) // only for directional light
        lvp := lp.Mul4(lv)
        shader.Matrix4f("light_vp", &lvp[0])

        /* set light uniform attributes */
        shader.Vec3("light.Position", &light.Position)
        shader.Vec3("light.Color",    &light.Color)
        shader.Int32("light.Type",     int32(light.Type))
        shader.Float("light.Range",    light.Range)
        shader.Float("light.attenuation.Constant",  light.Attenuation.Constant)
        shader.Float("light.attenuation.Linear",    light.Attenuation.Linear)
        shader.Float("light.attenuation.Quadratic", light.Attenuation.Quadratic)

        /* render light */
        gl.Viewport(0, 0, int32(scene.Camera.Width), int32(scene.Camera.Height))
        p.quad.Draw()
    }

    /* reset GL state */
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}