package engine

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/render"
)

type LineDrawable interface {
	object.Component
	DrawLines(render.Args) error
}

// LinePass draws line geometry
type LinePass struct {
}

// NewLinePass sets up a line geometry pass.
func NewLinePass() *LinePass {
	return &LinePass{}
}

// DrawPass executes the line pass
func (p *LinePass) Draw(args render.Args, scene object.T) {
	render.BindScreenBuffer()
	render.SetViewport(render.Viewport{
		Width:  args.Viewport.Width,
		Height: args.Viewport.Height,
	})

	objects := query.New[LineDrawable]().Collect(scene)
	for _, drawable := range objects {
		drawable.DrawLines(args.Apply(drawable.Object().Transform().World()))
	}
}
