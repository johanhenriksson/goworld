package engine

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/johanhenriksson/goworld/render"
)

// RenderPass is a step in the rendering pipeline
type RenderPass interface {
	DrawPass(*Scene)
}

// PassMap maps names to Render Passes
type PassMap map[string]RenderPass

// Renderer holds references to the Scene and is responsible for executing render passes in order
type Renderer struct {
	Passes  []RenderPass
	passMap PassMap
}

// NewRenderer instantiates a new rendering pipeline.
func NewRenderer() *Renderer {
	r := &Renderer{
		Passes:  []RenderPass{},
		passMap: make(PassMap),
	}
	return r
}

// Append a new render pass
func (r *Renderer) Append(name string, pass RenderPass) {
	if len(name) == 0 {
		panic(fmt.Errorf("Render passes must have names"))
	}

	r.Passes = append(r.Passes, pass)
	if len(name) > 0 {
		r.passMap[name] = pass
	}
}

// Get render pass by name
func (r *Renderer) Get(name string) RenderPass {
	return r.passMap[name]
}

// Reset the render pipeline.
func (r *Renderer) Reset() {
	r.Passes = []RenderPass{}
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

	for _, pass := range r.Passes {
		pass.DrawPass(scene)
	}
}
