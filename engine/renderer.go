package engine

import (
	"github.com/go-gl/gl/v4.1-core/gl"
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
	Scene   *Scene
	passMap PassMap
}

// NewRenderer instantiates a new rendering pipeline.
func NewRenderer(scene *Scene) *Renderer {
	r := &Renderer{
		Scene:   scene,
		Passes:  []RenderPass{},
		passMap: make(PassMap),
	}
	return r
}

// Append a new render pass
func (r *Renderer) Append(name string, pass RenderPass) {
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
func (r *Renderer) Draw() {
	/* Clear screen */
	gl.ClearColor(0.9, 0.9, 0.9, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	/* Enable blending */
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	/* Depth test */
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	for _, pass := range r.Passes {
		pass.DrawPass(r.Scene)
	}
}

// Update the world.
func (r *Renderer) Update(dt float32) {
	r.Scene.Update(dt)
}
