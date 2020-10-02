package engine

import (
	"github.com/johanhenriksson/goworld/render"

	"github.com/go-gl/gl/v4.1-core/gl"
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
	gl.Enable(gl.BLEND)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	// draw scene
	scene.DrawPass(render.ForwardPass)

	// reset
	// gl.Disable(gl.CULL_FACE)
}
