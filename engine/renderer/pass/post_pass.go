package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type PostPass interface {
	Pass
}

type postPass struct {
	target vulkan.Target
	prev   Pass
}

func NewPostPass(target vulkan.Target, prev Pass) PostPass {
	return &postPass{
		target: target,
		prev:   prev,
	}
}

func (p *postPass) Draw(args render.Args, scene object.T) {
	worker := p.target.Worker(args.Context.Index)
	worker.Submit(command.SubmitInfo{
		Marker: "GuiPass",
		Signal: []sync.Semaphore{args.Context.RenderComplete},
		Wait: []command.Wait{
			{
				Semaphore: p.prev.Completed(),
				Mask:      vk.PipelineStageFragmentShaderBit,
			},
		},
		Then: func() {
			args.Context.InFlight.Unlock()
		},
	})
}

func (p *postPass) Name() string {
	return "Pre"
}

func (p *postPass) Completed() sync.Semaphore {
	return nil
}

func (p *postPass) Destroy() {
}
