package engine

import (
	"github.com/johanhenriksson/goworld/render"
)

type LineDrawable interface {
	DrawLines(DrawArgs)
}

// LinePass draws line geometry
type LinePass struct {
	queue *DrawQueue
}

// NewLinePass sets up a line geometry pass.
func NewLinePass() *LinePass {
	return &LinePass{
		queue: NewDrawQueue(),
	}
}

func (p *LinePass) Type() render.Pass {
	return render.Line
}

func (p *LinePass) Resize(width, height int) {}

// DrawPass executes the line pass
func (p *LinePass) Draw(scene *Scene) {
	scene.Camera.Use()

	p.queue.Clear()
	scene.Collect(p)
	for _, cmd := range p.queue.items {
		drawable := cmd.Component.(LineDrawable)
		drawable.DrawLines(cmd.Args)
	}
}

func (p *LinePass) Visible(c Component, args DrawArgs) bool {
	_, ok := c.(LineDrawable)
	return ok
}

func (p *LinePass) Queue(c Component, args DrawArgs) {
	p.queue.Add(c, args)
}
