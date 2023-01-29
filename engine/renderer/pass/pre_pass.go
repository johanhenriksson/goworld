package pass

import (
	"github.com/johanhenriksson/goworld/core/object"
	"github.com/johanhenriksson/goworld/render"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/render/vulkan"

	vk "github.com/vulkan-go/vulkan"
)

type PrePass interface {
	Pass
}

type prePass struct {
	target    vulkan.Target
	completed sync.Semaphore
}

func NewPrePass(target vulkan.Target) PrePass {
	return &prePass{
		target:    target,
		completed: sync.NewSemaphore(target.Device()),
	}
}

func (p *prePass) Record(cmds command.Recorder, args render.Args, scene object.T) {

}

func (p *prePass) Draw(args render.Args, scene object.T) {
	var waits []command.Wait
	if args.Context.ImageAvailable != nil {
		waits = []command.Wait{
			{
				Semaphore: args.Context.ImageAvailable,
				Mask:      vk.PipelineStageColorAttachmentOutputBit,
			},
		}
	}

	worker := p.target.Worker(args.Context.Index)
	worker.Submit(command.SubmitInfo{
		Marker: "PrePass",
		Signal: []sync.Semaphore{p.completed},
		Wait:   waits,
	})
}

func (p *prePass) Name() string {
	return "Pre"
}

func (p *prePass) Completed() sync.Semaphore {
	return p.completed
}

func (p *prePass) Destroy() {
	p.completed.Destroy()
}
