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
func NewGeometryPass(bufferWidth, bufferHeight int32) *GeometryPass {
	p := &GeometryPass{
		Buffer: render.CreateGeometryBuffer(bufferWidth, bufferHeight),
	}
	return p
}

// DrawPass executes the geometry pass
func (p *GeometryPass) DrawPass(scene *Scene) {
	p.Buffer.Bind()
	p.Buffer.ClearColor = scene.Camera.Clear
	p.Buffer.Clear()

	// kind-of hack to clear the diffuse buffer separately
	// allows us to clear with the camera background color
	// other buffers need to be zeroed. or???
	/*
		camera := scene.Camera
		gl.DrawBuffer(gl.COLOR_ATTACHMENT0) // use only diffuse buffer
		gl.ClearColor(camera.Clear.R, camera.Clear.G, camera.Clear.B, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		p.Buffer.DrawBuffers()
	*/

	// setup rendering
	gl.Disable(gl.BLEND)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	// draw scene
	scene.DrawPass(render.GeometryPass)

	// reset
	gl.Disable(gl.CULL_FACE)
	gl.Enable(gl.BLEND)

	p.Buffer.Unbind()
}
