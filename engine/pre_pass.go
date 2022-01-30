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
	objects := object.NewQuery().
		Where(RequiresPreDraw).
		Collect(scene)

	for _, component := range objects {
		drawable := component.(PreDrawable)
		drawable.PreDraw(args.Apply(component.Object().Transform().World()))
	}
}

func RequiresPreDraw(c object.Component) bool {
	_, ok := c.(PreDrawable)
	return ok
}
