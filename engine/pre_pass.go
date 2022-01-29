package engine

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/scene"
	"github.com/johanhenriksson/goworld/render"
)

type PrePass struct{}

type PreDrawable interface {
	PreDraw(render.Args)
}

func (p *PrePass) Draw(args render.Args, scene scene.T) {
	query := object.NewQuery(func(c object.Component) bool {
		_, ok := c.(PreDrawable)
		return ok
	})
	scene.Collect(&query)

	args = ArgsWithCamera(args, scene.Camera())
	for _, component := range query.Results {
		drawable := component.(PreDrawable)
		drawable.PreDraw(args.Apply(component.Object().Transform().World()))
	}
}
