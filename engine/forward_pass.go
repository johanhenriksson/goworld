package engine

import (
	"github.com/johanhenriksson/goworld/render"
)

// ForwardPass holds information required to perform a forward rendering pass.
type ForwardPass struct{}

// NewForwardPass sets up a forward pass.
func NewForwardPass() *ForwardPass {
	return &ForwardPass{}
}

// DrawPass executes the forward pass
func (p *ForwardPass) DrawPass(scene *Scene) {
	scene.Camera.Use()
	render.ScreenBuffer.Bind()

	// setup rendering
	render.Blend(true)
	render.CullFace(render.CullBack)

	// draw scene
	scene.DrawPass(DrawForward)
}
