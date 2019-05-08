package engine

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render"
	//"github.com/go-gl/gl/v4.1-core/gl"
)

type LinePass struct {
	Material *render.Material
}

/* Sets up a geometry pass.
 * A geometry buffer of the given bufferWidth x bufferHeight will be created automatically */
func NewLinePass() *LinePass {
	mat := assets.GetMaterialCached("lines")
	p := &LinePass{
		Material: mat,
	}
	return p
}

func (p *LinePass) DrawPass(scene *Scene) {
	/* Draw scene */
	p.Material.Use()
	scene.Draw("lines", p.Material.Shader)
}
