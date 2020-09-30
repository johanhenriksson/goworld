package engine

import (
	// "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/johanhenriksson/goworld/render"
)

// LinePass draws line geometry
type LinePass struct{}

// NewLinePass sets up a line geometry pass.
func NewLinePass() *LinePass {
	return &LinePass{}
}

// DrawPass executes the line pass
func (p *LinePass) DrawPass(scene *Scene) {
	scene.Camera.Use()
	scene.DrawPass(render.LinePass)
}
