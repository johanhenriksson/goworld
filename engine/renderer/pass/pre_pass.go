package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/render"
)

type PrePass struct{}

type PreDrawable interface {
	object.Component
	PreDraw(render.Args, object.T) error
}

func (p *PrePass) Draw(args render.Args, scene object.T) {
	objects := query.New[PreDrawable]().Collect(scene)
	for _, component := range objects {
		component.PreDraw(args.Apply(component.Object().Transform().World()), scene)
	}
}
