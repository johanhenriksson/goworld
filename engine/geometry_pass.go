package engine

import (
	"github.com/johanhenriksson/goworld/render"

	"github.com/go-gl/gl/v4.1-core/gl"
)

// GeometryPass draws the scene geometry to a G-buffer
type GeometryPass struct {
	Buffer *render.GeometryBuffer
}

// NewGeometryPass sets up a geometry pass.
func NewGeometryPass(bufferWidth, bufferHeight int) *GeometryPass {
	p := &GeometryPass{
		Buffer: render.CreateGeometryBuffer(bufferWidth, bufferHeight),
	}
	return p
}

// DrawPass executes the geometry pass
func (p *GeometryPass) DrawPass(scene *Scene) {
	p.Buffer.Bind()
	render.Clear()
	render.ClearDepth()

	// kind-of hack to clear the diffuse buffer separately
	// allows us to clear with the camera background color
	// other buffers need to be zeroed. or???
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0) // use only diffuse buffer
	render.ClearWith(scene.Camera.Clear)

	p.Buffer.DrawBuffers()

	// setup rendering
	render.Blend(false)
	render.CullFace(render.CullBack)
	render.DepthOutput(true)

	// draw scene
	scene.DrawPass(DrawGeometry)

	p.Buffer.Unbind()
}
