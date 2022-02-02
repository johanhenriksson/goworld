package engine

import (
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/engine/screen_quad"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/backend/gl/gl_shader"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// OutputPass is the final pass that writes to a camera frame buffer.
type OutputPass struct {
	Color  texture.T
	Depth  texture.T
	shader shader.T
	quad   screen_quad.T
	mat    material.T
}

// NewOutputPass creates a new output pass for the given input texture.
// The output pass writes a full screen texture to the screen, and restores the depth buffer.
func NewOutputPass(color texture.T, depth texture.T) *OutputPass {
	shader := gl_shader.CompileShader(
		"output_pass",
		"/assets/shaders/pass/postprocess.vs",
		"/assets/shaders/pass/output.fs")

	mat := material.New("output_pass", shader)
	mat.Texture("tex_input", color)
	mat.Texture("tex_depth", depth)

	return &OutputPass{
		Color:  color,
		Depth:  depth,
		shader: shader,
		mat:    mat,
		quad:   screen_quad.New(shader),
	}
}

// DrawPass draws the input texture to the scene camera buffer.
func (p *OutputPass) Draw(args render.Args, scene scene.T) {
	render.BindScreenBuffer()
	render.SetViewport(0, 0, args.Viewport.FrameWidth, args.Viewport.FrameHeight)

	// ensures we dont fail depth tests while restoring the depth buffer
	gl.DepthFunc(gl.ALWAYS)

	// draw
	p.mat.Use()
	p.quad.Draw()

	gl.DepthFunc(gl.LESS)
}
