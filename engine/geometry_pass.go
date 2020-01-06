package engine

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type GeometryPass struct {
	Buffer   *render.GeometryBuffer
	Material *render.Material
}

/* Sets up a geometry pass.
 * A geometry buffer of the given bufferWidth x bufferHeight will be created automatically */
func NewGeometryPass(bufferWidth, bufferHeight int32) *GeometryPass {
	mat := assets.GetMaterial("ssao_color_geometry")
	p := &GeometryPass{
		Buffer:   render.CreateGeometryBuffer(bufferWidth, bufferHeight),
		Material: mat,
	}
	return p
}

func (p *GeometryPass) DrawPass(scene *Scene) {
	p.Buffer.Bind()
	//p.Buffer.ClearColor = scene.Camera.Clear
	p.Buffer.Clear()

	// kind-of hack to clear the diffuse buffer separately
	// why???
	camera := scene.Camera
	gl.DrawBuffer(gl.COLOR_ATTACHMENT0) //
	gl.ClearColor(camera.Clear.R, camera.Clear.G, camera.Clear.B, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	p.Buffer.DrawBuffers()

	//p.Material.Use()

	// setup rendering
	gl.Disable(gl.BLEND)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	// draw scene
	scene.Draw(render.GeometryPass, p.Material.ShaderProgram)

	// reset
	gl.Disable(gl.CULL_FACE)
	gl.Enable(gl.BLEND)

	p.Buffer.Unbind()
}
