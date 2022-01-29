package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/random"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	glshader "github.com/johanhenriksson/goworld/render/backend/gl/shader"
	gltex "github.com/johanhenriksson/goworld/render/backend/gl/texture"
	"github.com/johanhenriksson/goworld/render/backend/types"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
)

// SSAOSettings holds parameters for SSAO.
type SSAOSettings struct {
	Samples int
	Scale   int
	Radius  float32
	Bias    float32
	Power   float32
}

// SSAOPass renders a screen space ambient occlusion texture for a gbuffer-based scene.
type SSAOPass struct {
	SSAOSettings

	GBuffer  *render.GeometryBuffer
	Gaussian *GaussianPass
	Output   texture.T
	Noise    texture.T
	Kernel   []vec3.T

	fbo    *render.FrameBuffer
	shader shader.T
	mat    material.T
	quad   *Quad
}

// NewSSAOPass creates a new SSAO pass from a gbuffer and SSAO settings.
func NewSSAOPass(gbuff *render.GeometryBuffer, settings *SSAOSettings) *SSAOPass {
	fbo := render.CreateFrameBuffer(gbuff.Width/settings.Scale, gbuff.Height/settings.Scale)
	texture := fbo.NewBuffer(gl.COLOR_ATTACHMENT0, gl.RED, gl.RGB, gl.FLOAT) // diffuse (rgb)

	// gaussian blur pass
	gaussian := NewGaussianPass(texture)

	// generate sample kernel
	kernel := createSSAOKernel(settings.Samples)

	// generate noise texture
	noise := createHemisphereNoiseTexture(4)

	shader := glshader.CompileShader(
		"ssao_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/ssao.fs")

	mat := material.New("ssao_pass", shader)
	mat.Texture("tex_normal", gbuff.Normal)
	mat.Texture("tex_position", gbuff.Position)
	mat.Texture("tex_noise", noise)

	// create a render quad

	p := &SSAOPass{
		SSAOSettings: *settings,
		GBuffer:      gbuff,
		quad:         NewQuad(shader),

		Noise:  noise,
		Kernel: kernel,
		Output: texture,

		Gaussian: gaussian,

		fbo:    fbo,
		shader: shader,
		mat:    mat,
	}

	return p
}

// DrawPass draws the SSAO texture.
func (p *SSAOPass) Draw(args render.Args, scene scene.T) {
	render.Blend(false)
	render.DepthOutput(false)

	// update projection
	p.mat.Use()
	if err := p.mat.Int32("kernel_size", len(p.Kernel)); err != nil {
		panic(err)
	}
	if err := p.mat.Float("bias", p.Bias); err != nil {
		panic(err)
	}
	p.mat.Float("radius", p.Radius)
	p.mat.Float("power", p.Power)
	if err := p.mat.Int32("scale", p.Scale); err != nil {
		panic(err)
	}
	if err := p.mat.Vec3Array("samples", p.Kernel); err != nil {
		panic(err)
	}
	p.mat.Mat4("projection", scene.Camera().Projection())

	// run occlusion pass
	p.fbo.Bind()
	defer p.fbo.Unbind()
	p.fbo.Resize(p.GBuffer.Width/p.Scale, p.GBuffer.Height/p.Scale)

	render.ClearWith(color.White)
	p.quad.Draw()

	// run blur pass
	p.Gaussian.DrawPass(args, scene)

	render.DepthOutput(true)
}

func createSSAOKernel(samples int) []vec3.T {
	kernel := make([]vec3.T, samples)
	for i := 0; i < len(kernel); i++ {
		sample := vec3.Random(vec3.New(-1, -1, 0), vec3.One)
		sample.Normalize()
		sample = sample.Scaled(random.Range(0, 1)) // random length

		// scaling
		scale := float32(i) / float32(samples)
		scale = math.Lerp(0.1, 1.0, scale*scale)
		sample = sample.Scaled(scale)

		kernel[i] = sample
	}
	return kernel
}

func createHemisphereNoiseTexture(size int) texture.T {
	noise := gltex.New(size, size)
	noise.SetInternalFormat(gl.RGB16F)
	noise.SetFormat(texture.RGB)
	noise.SetDataType(types.Float)

	noise.SetFilter(texture.NearestFilter)
	noise.SetWrapMode(texture.RepeatWrap)

	noiseData := make([]float32, 3*size*size)
	for i := 0; i < len(noiseData); i += 3 {
		noiseData[i+0] = random.Range(-1, 1)
		noiseData[i+1] = 1 // random.Range(-1, 1)
		noiseData[i+2] = 0
	}
	noise.BufferFloats(noiseData)

	return noise
}
