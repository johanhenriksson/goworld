package engine

import (
	"fmt"
	"math/rand"
	"github.com/go-gl/gl/v4.1-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/johanhenriksson/goworld/render"
)

type SSAOPass struct {
	Output *render.FrameBuffer
	Texture *render.Texture
	Material *render.Material
	Quad *render.RenderQuad
	Noise *render.Texture
	Kernel []mgl.Vec3
	Radius float32
	Bias float32
}

func NewSSAOPass(gbuff *render.GeometryBuffer) *SSAOPass {
	ssaoFbo := render.CreateFrameBuffer(gbuff.Width, gbuff.Height)
	ssaoFbo.ClearColor = render.Color4(0, 0, 0, 1)
	ssaoBuffer := ssaoFbo.AddBuffer(gl.COLOR_ATTACHMENT0, gl.RGB, gl.RGB, gl.FLOAT) // diffuse (rgb)

	/* use a virtual material to help with vertex attributes and textures */
	mat := render.CreateMaterial(render.CompileVFShader("/assets/shaders/ssao"))

	/* we're going to render a simple quad, so we input
	 * position and texture coordinates */
	mat.AddDescriptor("position", gl.FLOAT, 3, 20, 0, false, false)
	mat.AddDescriptor("texcoord", gl.FLOAT, 2, 20, 12, false, false)

	/* the shader uses 3 textures from the geometry frame buffer.
	 * they are previously rendered in the geometry pass. */
	mat.AddTexture("tex_position", gbuff.Position)
	mat.AddTexture("tex_normal", gbuff.Normal)

	// sample kernel
	kernel := make([]mgl.Vec3, 64)
	for i := 0; i < len(kernel); i++ {
		sample := mgl.Vec3{
			rand.Float32() * 2 - 1,
			rand.Float32() * 2 - 1,
			rand.Float32(),
		}
		sample = sample.Normalize()
		sample.Mul(rand.Float32()) // random length

		// scaling
		scale := float32(i) / 64.0
		scale = lerp(0.1, 1.0, scale * scale)
		sample = sample.Mul(scale)
		
		kernel[i] = sample
	}

	// noise texture
	nsize := int32(8)
	noise := render.CreateTexture(nsize, nsize)
	noise.InternalFormat = gl.RGB32F
	noise.Format = gl.RGB
	noise.DataType = gl.FLOAT
	noiseData := make([]float32, 3 * nsize * nsize)
	for i := 0; i < len(noiseData); i += 3 {
		noiseData[i+0] = rand.Float32() * 2 - 1
		noiseData[i+1] = rand.Float32() * 2 - 1
		noiseData[i+2] = 0
	}
	noise.BufferFloats(noiseData)

	mat.AddTexture("tex_noise", noise)

	/* create a render quad */
	quad := render.NewRenderQuad()
	/* set up vertex attribute pointers */
	mat.SetupVertexPointers()

	return &SSAOPass{
		Output: ssaoFbo,
		Texture: ssaoBuffer,
		Material: mat,
		Quad: quad,
		Noise: noise,
		Kernel: kernel,
		Radius: 3.0,
		Bias: 0.03,
	}
}

func (p *SSAOPass) DrawPass(scene *Scene) {
	p.Output.Bind()
	p.Output.Clear()

	p.Material.Use()
	shader := p.Material.Shader

	shader.Matrix4f("projection", &scene.Camera.Projection[0])
	shader.Int32("kernel_size", int32(len(p.Kernel)))
	shader.Float("bias", p.Bias)
	shader.Float("radius", p.Radius)

	for i := 0; i < 64; i++ {
		shader.Vec3(fmt.Sprintf("samples[%d]", i), &p.Kernel[i])
	}

	p.Quad.Draw()

	p.Output.Unbind()
}

func lerp(a, b, f float32) float32 {
    return a + f * (b - a)
}