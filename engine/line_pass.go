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

// DrawPass executes the line pass
func (p *LinePass) Draw(args render.Args, scene scene.T) {
	render.BindScreenBuffer()
	render.SetViewport(0, 0, args.Viewport.FrameWidth, args.Viewport.FrameHeight)

	objects := object.NewQuery().
		Where(IsLineDrawable).
		Collect(scene)

	for _, component := range objects {
		drawable := component.(LineDrawable)
		drawable.DrawLines(args.Apply(component.Object().Transform().World()))
	}
}

func IsLineDrawable(c object.Component) bool {
	_, ok := c.(LineDrawable)
	return ok
}
