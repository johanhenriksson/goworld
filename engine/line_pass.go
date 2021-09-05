package engine

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/render"
)

type LineDrawable interface {
	DrawLines(render.Args)
}

// LinePass draws line geometry
type LinePass struct {
}

// NewLinePass sets up a line geometry pass.
func NewLinePass() *LinePass {
	return &LinePass{}
}

func (p *LinePass) Resize(width, height int) {}

// DrawPass executes the line pass
func (p *LinePass) Draw(scene scene.T) {
	// scene.Camera.Use()
	render.ScreenBuffer.Bind()

	query := object.NewQuery(func(c object.Component) bool {
		_, ok := c.(LineDrawable)
		return ok
	})
	scene.Collect(&query)

	args := ArgsFromCamera(scene.Camera())
	for _, component := range query.Results {
		drawable := component.(LineDrawable)
		drawable.DrawLines(args.Apply(component.Object().Transform().World()))
	}
}
