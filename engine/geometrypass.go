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
	p.Buffer.Clear()

	p.Material.Use()

	/* Draw scene */
	gl.Disable(gl.BLEND)
	scene.Draw("geometry", p.Material.Shader)
	gl.Enable(gl.BLEND)

	p.Buffer.Unbind()
}
