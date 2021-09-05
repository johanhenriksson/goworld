package engine

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/render"
)

// PassMap maps names to Render Passes
type PassMap map[string]DrawPass

type PostPass interface {
	Input(*render.Texture)
	Output() *render.Texture
}

// Renderer holds references to the Scene and is responsible for executing render passes in order
type Renderer struct {
	Passes  []DrawPass
	passMap PassMap

	Output    *OutputPass
	Geometry  *GeometryPass
	Light     *LightPass
	Forward   *ForwardPass
	SSAO      *SSAOPass
	Particles *ParticlePass
	Colors    *ColorPass
	Lines     *LinePass
}

// NewRenderer instantiates a new rendering pipeline.
func NewRenderer() *Renderer {
	r := &Renderer{
		Passes:  []DrawPass{},
		passMap: make(PassMap),
	}

	width, height := render.ScreenBuffer.Width, render.ScreenBuffer.Height

	r.Geometry = NewGeometryPass(width, height)
	r.Light = NewLightPass(r.Geometry.Buffer)

	r.Forward = NewForwardPass(r.Geometry.Buffer, r.Light.Output)

	r.SSAO = NewSSAOPass(r.Geometry.Buffer, &SSAOSettings{
		Samples: 16,
		Radius:  0.1,
		Bias:    0.03,
		Power:   2.0,
		Scale:   2,
	})

	r.Colors = NewColorPass(r.Light.Output, "saturated", r.SSAO.Gaussian.Output)
	r.Output = NewOutputPass(r.Colors.Output, r.Geometry.Buffer)

	r.Lines = NewLinePass()
	// r.Particles = NewParticlePass()

	return r
}

func (r *Renderer) Resize(width, height int) {
	r.Geometry.Resize(width, height)
	r.Light.Resize(width, height)
	r.Colors.Resize(width, height)
}

// Append a new render pass
func (r *Renderer) Append(name string, pass DrawPass) {
	if len(name) == 0 {
		panic(fmt.Errorf("render passes must have names"))
	}

	r.Passes = append(r.Passes, pass)
	if len(name) > 0 {
		r.passMap[name] = pass
	}
}

// Get render pass by name
func (r *Renderer) Get(name string) DrawPass {
	return r.passMap[name]
}

// Reset the render pipeline.
func (r *Renderer) Reset() {
	r.Passes = []DrawPass{}
	r.passMap = make(PassMap)
}

// Draw the world.
func (r *Renderer) Draw(scene *Scene) {
	// clear screen buffer
	render.ScreenBuffer.Bind()
	gl.ClearColor(0.9, 0.9, 0.9, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// enable blending
	render.Blend(true)
	render.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// enable depth test
	render.DepthTest(true)
	gl.DepthFunc(gl.LESS)

	r.Geometry.Draw(scene)
	r.Light.Draw(scene)
	r.Forward.Draw(scene)
	r.SSAO.Draw(scene)
	r.Colors.Draw(scene)
	r.Output.Draw(scene)
	r.Lines.Draw(scene)
	// r.Particles.Draw(scene)

	for _, pass := range r.Passes {
		pass.Draw(scene)
	}
}
