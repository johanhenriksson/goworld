package deferred

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/core/camera"
	"github.com/johanhenriksson/goworld/core/light"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_shader"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/framebuffer"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/texture"
)

type LightShader interface {
	material.T
	SetCamera(camera.T)
	SetLightDescriptor(light.Descriptor)
	SetShadowMap(texture.T)
	SetShadowStrength(float32)
	SetShadowBias(float32)
}

type lightshader struct {
	material.T
}

func NewLightShader(gbuffer framebuffer.Geometry) LightShader {
	shader := gl_shader.CompileShader(
		"light_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/light.fs")

	mat := material.New("light_pass", shader)
	mat.Texture("tex_diffuse", gbuffer.Diffuse())
	mat.Texture("tex_normal", gbuffer.Normal())
	mat.Texture("tex_depth", gbuffer.Depth())
	mat.Texture("tex_shadow", assets.GetColorTexture(color.White))

	return &lightshader{
		T: mat,
	}
}

func (sh *lightshader) SetShadowStrength(strength float32) {
	sh.Float("shadow_strength", strength)
}

func (sh *lightshader) SetShadowBias(bias float32) {
	sh.Float("shadow_bias", bias)
}

func (sh *lightshader) SetCamera(cam camera.T) {
	// why is this not equal to args.View/args.VP
	vInv := cam.ViewInv()
	vpInv := cam.ViewProjInv()
	sh.Mat4("cameraInverse", vpInv)
	sh.Mat4("viewInverse", vInv)
}

func (sh *lightshader) SetShadowMap(tex texture.T) {
	sh.Texture("tex_shadow", tex)
}

func (sh *lightshader) SetLightDescriptor(desc light.Descriptor) {
	sh.Int32("light.Type", int(desc.Type))
	sh.Vec3("light.Position", desc.Position)
	sh.RGB("light.Color", desc.Color)
	sh.Float("light.Intensity", desc.Intensity)

	if desc.Type == light.Directional {
		sh.Mat4("light_vp", desc.ViewProj)
	}

	if desc.Type == light.Point {
		sh.Float("light.Range", desc.Range)
		sh.Float("light.attenuation.Constant", desc.Attenuation.Constant)
		sh.Float("light.attenuation.Linear", desc.Attenuation.Linear)
		sh.Float("light.attenuation.Quadratic", desc.Attenuation.Quadratic)
	}
}
