package engine

import (
	"fmt"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/math/vec3"
	"github.com/johanhenriksson/goworld/render"
)

// SSAOSettings holds parameters for SSAO.
type SSAOSettings struct {
	Samples int
	Scale   int32
	Radius  float32
	Bias    float32
	Power   float32
}

// SSAOPass renders a screen space ambient occlusion texture for a gbuffer-based scene.
type SSAOPass struct {
	SSAOSettings

	fbo *render.FrameBuffer

	Output   *render.Texture
	Material *render.Material
	Quad     *render.Quad
	Noise    *render.Texture
	Kernel   []vec3.T

	Gaussian *GaussianPass
}

// NewSSAOPass creates a new SSAO pass from a gbuffer and SSAO settings.
func NewSSAOPass(gbuff *render.GeometryBuffer, settings *SSAOSettings) *SSAOPass {
	fbo := render.CreateFrameBuffer(gbuff.Width/settings.Scale, gbuff.Height/settings.Scale)
	fbo.ClearColor = render.Color4(1, 1, 1, 1)
	texture := fbo.AttachBuffer(gl.COLOR_ATTACHMENT0, gl.RED, gl.RGB, gl.FLOAT) // diffuse (rgb)

	// gaussian blur pass
	gaussian := NewGaussianPass(texture)

	// generate sample kernel
	kernel := createSSAOKernel(settings.Samples)

	// generate noise texture
	noise := createHemisphereNoiseTexture(4)

	/* use a virtual material to help with vertex attributes and textures */
	mat := render.CreateMaterial(render.CompileVFShader("/assets/shaders/ssao"))
	mat.AddDescriptors(render.F32_XYZUV)
	mat.AddTexture("tex_position", gbuff.Position)
	mat.AddTexture("tex_normal", gbuff.Normal)
	mat.AddTexture("tex_noise", noise)

	/* create a render quad */
	quad := render.NewQuad(mat)

	p := &SSAOPass{
		SSAOSettings: *settings,

		fbo: fbo,

		Material: mat,
		Quad:     quad,
		Noise:    noise,
		Kernel:   kernel,
		Output:   texture,

		Gaussian: gaussian,
	}

	// set up shader uniforms
	mat.Use()
	mat.Int32("kernel_size", int32(len(p.Kernel)))
	mat.Float("bias", p.Bias)
	mat.Float("radius", p.Radius)
	mat.Float("power", p.Power)
	mat.Int32("scale", p.Scale)

	for i := 0; i < len(p.Kernel); i++ {
		mat.Vec3(fmt.Sprintf("samples[%d]", i), &p.Kernel[i])
	}

	return p
}

// DrawPass draws the SSAO texture.
func (p *SSAOPass) DrawPass(scene *Scene) {
	// update projection
	p.Material.Use()
	p.Material.Mat4f("projection", &scene.Camera.Projection)

	// run occlusion pass
	p.fbo.Bind()
	p.fbo.Clear()
	p.Quad.Draw()
	p.fbo.Unbind()

	// run blur pass
	p.Gaussian.DrawPass(scene)
}

func createSSAOKernel(samples int) []vec3.T {
	kernel := make([]vec3.T, samples)
	for i := 0; i < len(kernel); i++ {
		sample := vec3.T{
			X: rand.Float32()*2 - 1,
			Y: rand.Float32()*2 - 1,
			Z: rand.Float32(),
		}
		sample.Normalize()
		sample = sample.Scaled(rand.Float32()) // random length

		// scaling
		scale := float32(i) / float32(samples)
		scale = lerp(0.1, 1.0, scale*scale)
		sample = sample.Scaled(scale)

		kernel[i] = sample
	}
	return kernel
}

func createHemisphereNoiseTexture(size int) *render.Texture {
	noise := render.CreateTexture(int32(size), int32(size))
	noise.InternalFormat = gl.RGB16F
	noise.Format = gl.RGB
	noise.DataType = gl.FLOAT

	noise.SetFilter(render.NearestFilter)
	noise.SetWrapMode(render.RepeatWrap)

	noiseData := make([]float32, 3*size*size)
	for i := 0; i < len(noiseData); i += 3 {
		noiseData[i+0] = rand.Float32()*2 - 1
		noiseData[i+1] = rand.Float32()*2 - 1
		noiseData[i+2] = 0
	}
	noise.BufferFloats(noiseData)

	return noise
}

func lerp(a, b, f float32) float32 {
	return a + f*(b-a)
}
