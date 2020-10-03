package engine

import (
	"github.com/johanhenriksson/goworld/render"
)

// GeometryPass draws the scene geometry to a G-buffer
type ForwardPass struct{}

// NewGeometryPass sets up a geometry pass.
func NewForwardPass() *ForwardPass {
	return &ForwardPass{}
}

// DrawPass executes the geometry pass
func (p *ForwardPass) DrawPass(scene *Scene) {
	scene.Camera.Use()

	// setup rendering
	render.Blend(true)
	render.CullFace(render.CullBack)

	// draw scene
	scene.DrawPass(render.ForwardPass)
}
