package engine

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/swapchain"
	"github.com/johanhenriksson/goworld/render/sync"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type TargetSize struct {
	Width  int
	Height int
	Frames int
	Scale  float32
}

type Target interface {
	Size() TargetSize
	Scale() float32
	Width() int
	Height() int
	Frames() int

	Surfaces() image.Array
	SurfaceFormat() core1_0.Format
	Aquire(command.Worker) (*swapchain.Context, error)
	Present(command.Worker, *swapchain.Context)

	Destroy()
}

type RenderTarget struct {
	device   *device.Device
	size     TargetSize
	format   core1_0.Format
	usage    core1_0.ImageUsageFlags
	surfaces image.Array
	contexts chan *swapchain.Context
}

var _ Target = (*RenderTarget)(nil)

func NewDepthTarget(device *device.Device, key string, size TargetSize) *RenderTarget {
	format := device.GetDepthFormat()
	usage := core1_0.ImageUsageSampled | core1_0.ImageUsageDepthStencilAttachment | core1_0.ImageUsageInputAttachment
	target, err := NewRenderTarget(device, key, format, usage, size)
	if err != nil {
		panic(err)
	}
	return target
}

func NewColorTarget(device *device.Device, key string, format core1_0.Format, size TargetSize) *RenderTarget {
	usage := core1_0.ImageUsageSampled | core1_0.ImageUsageColorAttachment | core1_0.ImageUsageInputAttachment | core1_0.ImageUsageTransferSrc
	target, err := NewRenderTarget(device, key, format, usage, size)
	if err != nil {
		panic(err)
	}
	return target
}

func NewRenderTarget(device *device.Device, key string, format core1_0.Format, usage core1_0.ImageUsageFlags, size TargetSize) (*RenderTarget, error) {
	var err error
	outputs := make(image.Array, size.Frames)
	contexts := make(chan *swapchain.Context, size.Frames) // context ring buffer channel
	for i := 0; i < size.Frames; i++ {
		outputs[i], err = image.New2D(device, fmt.Sprintf("%s:%d", key, i), size.Width, size.Height, format, false, usage)
		if err != nil {
			return nil, err
		}
		// send the context to the ring buffer
		// guaranteed to be non-blocking since the buffer is the same size as the number of frames
		contexts <- swapchain.NewContext(device, i)
	}

	return &RenderTarget{
		device:   device,
		size:     size,
		format:   format,
		usage:    usage,
		surfaces: outputs,
		contexts: contexts,
	}, nil
}

func (r *RenderTarget) Frames() int    { return len(r.surfaces) }
func (r *RenderTarget) Width() int     { return r.size.Width }
func (r *RenderTarget) Height() int    { return r.size.Height }
func (r *RenderTarget) Scale() float32 { return r.size.Scale }

func (r *RenderTarget) Size() TargetSize {
	return r.size
}

func (r *RenderTarget) Destroy() {
	// wait for each context to complete, then destroy it
	for i := 0; i < r.Frames(); i++ {
		ctx := <-r.contexts
		ctx.Destroy()
	}
	close(r.contexts)

	for _, output := range r.surfaces {
		output.Destroy()
	}
	r.surfaces = nil
}

func (i *RenderTarget) Surfaces() image.Array         { return i.surfaces }
func (i *RenderTarget) SurfaceFormat() core1_0.Format { return i.format }

func (i *RenderTarget) Aquire(worker command.Worker) (*swapchain.Context, error) {
	// wait for the next context to be available
	context := <-i.contexts

	// submit a command to signal the image available semaphore, mimicing the behavior of a swapchain
	worker.Submit(command.SubmitInfo{
		Marker:   "AquireRenderTarget",
		Commands: command.Empty,
		Signal: []*sync.Semaphore{
			context.ImageAvailable,
		},
	})

	return context, nil
}

func (t *RenderTarget) Present(worker command.Worker, context *swapchain.Context) {
	// wait for the render complete semaphore to be signaled
	// then, return the context to the ring buffer
	worker.Submit(command.SubmitInfo{
		Marker:   "PresentRenderTarget",
		Commands: command.Empty,
		Wait: []command.Wait{
			{
				Semaphore: context.RenderComplete,
				Mask:      core1_0.PipelineStageColorAttachmentOutput,
			},
		},
		Callback: func() {
			// guaranteed to be non-blocking since the buffer is the same size as the number of frames
			t.contexts <- context
		},
	})
}
