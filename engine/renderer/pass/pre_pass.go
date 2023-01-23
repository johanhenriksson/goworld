package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/core/object/query"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/vulkan"
)

type PrePass interface {
	Pass
}

type prePass struct {
	target vulkan.Target
}

type PreDrawable interface {
	object.Component
	PreDraw(render.Args, object.T) error
}

func NewPrePass(target vulkan.Target) PrePass {
	return &prePass{
		target: target,
	}
}

func (p *prePass) Draw(args render.Args, scene object.T) {
	objects := query.New[PreDrawable]().Collect(scene)
	for _, component := range objects {
		component.PreDraw(args.Apply(component.Object().Transform().World()), scene)
	}
}

func (p *prePass) Name() string {
	return "Pre"
}

func (p *prePass) Completed() sync.Semaphore {
	return nil
}

func (p *prePass) Destroy() {
}
