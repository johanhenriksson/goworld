package engine

import (
	"fmt"
	"math/rand"

	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/render"
)

type SSAOSettings struct {
	Samples int
	Scale   int32
	Radius  float32
	Bias    float32
	Power   float32
}

type SSAOPass struct {
	SSAOSettings

	fbo *render.FrameBuffer

	Output   *render.Texture
	Material *render.Material
	Quad     *render.Quad
	Noise    *render.Texture
	Kernel   []mgl.Vec3

	Gaussian *GaussianPass
}

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

func (p *SSAOPass) DrawPass(scene *Scene) {
	// update projection
	p.Material.Use()
	p.Material.Mat4f("projection", scene.Camera.Projection)

	// run occlusion pass
	p.fbo.Bind()
	p.fbo.Clear()
	p.Quad.Draw()
	p.fbo.Unbind()

	// run blur pass
	p.Gaussian.DrawPass(scene)
}

func createSSAOKernel(samples int) []mgl.Vec3 {
	kernel := make([]mgl.Vec3, samples)
	for i := 0; i < len(kernel); i++ {
		sample := mgl.Vec3{
			rand.Float32()*2 - 1,
			rand.Float32()*2 - 1,
			rand.Float32(),
		}
		sample = sample.Normalize()
		sample.Mul(rand.Float32()) // random length

		// scaling
		scale := float32(i) / float32(samples)
		scale = lerp(0.1, 1.0, scale*scale)
		sample = sample.Mul(scale)

		kernel[i] = sample
	}
	return kernel
}

func createHemisphereNoiseTexture(size int) *render.Texture {
	noise := render.CreateTexture(int32(size), int32(size))
	noise.InternalFormat = gl.RGB16F
	noise.Format = gl.RGB
	noise.DataType = gl.FLOAT

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)

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
