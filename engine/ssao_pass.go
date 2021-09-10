package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/math"
	"github.com/johanhenriksson/goworld/math/random"
	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/color"
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
	Output   *render.Texture
	Noise    *render.Texture
	Kernel   []vec3.T

	fbo      *render.FrameBuffer
	shader   *render.Shader
	textures *render.TextureMap
	quad     *Quad
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

	shader := render.CompileShader(
		"ssao_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/ssao.fs")

	tx := render.NewTextureMap(shader)
	tx.Add("tex_normal", gbuff.Normal)
	tx.Add("tex_position", gbuff.Position)
	tx.Add("tex_noise", noise)

	// create a render quad

	p := &SSAOPass{
		SSAOSettings: *settings,
		GBuffer:      gbuff,
		quad:         NewQuad(shader),

		Noise:  noise,
		Kernel: kernel,
		Output: texture,

		Gaussian: gaussian,

		fbo:      fbo,
		shader:   shader,
		textures: tx,
	}

	// set up shader uniforms
	shader.Use()
	shader.Int32("kernel_size", len(p.Kernel))
	shader.Float("bias", p.Bias)
	shader.Float("radius", p.Radius)
	shader.Float("power", p.Power)
	shader.Int32("scale", p.Scale)
	shader.Vec3Array("samples", p.Kernel)

	return p
}

// DrawPass draws the SSAO texture.
func (p *SSAOPass) Draw(scene scene.T) {
	render.Blend(false)
	render.DepthOutput(false)

	// update projection
	p.shader.Use()
	proj := scene.Camera().Projection()
	p.shader.Mat4("projection", &proj)
	p.textures.Use()

	// run occlusion pass
	p.fbo.Bind()

	render.ClearWith(color.White)
	p.quad.Draw()

	p.fbo.Unbind()

	// run blur pass
	p.Gaussian.DrawPass(scene)

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

func createHemisphereNoiseTexture(size int) *render.Texture {
	noise := render.CreateTexture(size, size)
	noise.InternalFormat = gl.RGB16F
	noise.Format = gl.RGB
	noise.DataType = gl.FLOAT

	noise.SetFilter(render.NearestFilter)
	noise.SetWrapMode(render.RepeatWrap)

	noiseData := make([]float32, 3*size*size)
	for i := 0; i < len(noiseData); i += 3 {
		noiseData[i+0] = random.Range(-1, 1)
		noiseData[i+1] = 1 // random.Range(-1, 1)
		noiseData[i+2] = 0
	}
	noise.BufferFloats(noiseData)

	return noise
}
