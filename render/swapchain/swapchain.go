package swapchain

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/util"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/khr_surface"
	"github.com/vkngwrapper/extensions/v2/khr_swapchain"
)

type T interface {
	device.Resource[khr_swapchain.Swapchain]

	Aquire(command.Worker) (*Context, error)
	Present(command.Worker, *Context)
	Resize(int, int)

	Images() image.Array
	SurfaceFormat() core1_0.Format
}

type swapchain struct {
	device     *device.Device
	ptr        khr_swapchain.Swapchain
	ext        khr_swapchain.Extension
	surface    khr_surface.Surface
	surfaceFmt khr_surface.SurfaceFormat
	images     image.Array
	frames     int
	width      int
	height     int
	resized    bool

	nextContext int
	contexts    []*Context
}

func New(device *device.Device, frames, width, height int, surface khr_surface.Surface, surfaceFormat khr_surface.SurfaceFormat) T {
	s := &swapchain{
		device:     device,
		ext:        khr_swapchain.CreateExtensionFromDevice(device.Ptr()),
		surface:    surface,
		surfaceFmt: surfaceFormat,
		frames:     frames,
		width:      width,
		height:     height,
	}
	s.create()
	return s
}

func (s *swapchain) Ptr() khr_swapchain.Swapchain {
	return s.ptr
}

func (s *swapchain) Images() image.Array           { return s.images }
func (s *swapchain) SurfaceFormat() core1_0.Format { return core1_0.Format(s.surfaceFmt.Format) }

func (s *swapchain) Resize(width, height int) {
	// resizing actually happens the next time a frame is aquired
	s.width = width
	s.height = height
	s.resized = true
}

func (s *swapchain) recreate() {
	log.Println("recreating swapchain")

	// recreate swapchain resources
	s.device.WaitIdle()
	s.Destroy()
	s.create()
}

func (s *swapchain) create() {
	imageFormat := core1_0.Format(s.surfaceFmt.Format)
	imageUsage := core1_0.ImageUsageColorAttachment | core1_0.ImageUsageTransferSrc
	imageSharing := core1_0.SharingModeExclusive

	swapInfo := khr_swapchain.SwapchainCreateInfo{
		Surface:         s.surface,
		MinImageCount:   s.frames,
		ImageFormat:     imageFormat,
		ImageColorSpace: khr_surface.ColorSpace(s.surfaceFmt.ColorSpace),
		ImageExtent: core1_0.Extent2D{
			Width:  s.width,
			Height: s.height,
		},
		ImageArrayLayers: 1,
		ImageUsage:       imageUsage,
		ImageSharingMode: imageSharing,
		PresentMode:      khr_surface.PresentModeFIFO,
		PreTransform:     khr_surface.TransformIdentity,
		CompositeAlpha:   khr_surface.CompositeAlphaOpaque,
		Clipped:          true,
	}

	var chain khr_swapchain.Swapchain
	chain, _, err := s.ext.CreateSwapchain(s.device.Ptr(), nil, swapInfo)
	if err != nil {
		panic(err)
	}
	s.ptr = chain
	s.resized = false

	swapimages, result, err := chain.SwapchainImages()
	if err != nil {
		panic(err)
	}
	if result != core1_0.VKSuccess {
		panic("failed to get swapchain images")
	}
	if len(swapimages) != s.frames {
		panic("failed to get the requested number of swapchain images")
	}

	// create images from swapchain buffers
	s.images = util.Map(swapimages, func(img core1_0.Image) *image.Image {
		return image.Wrap(s.device, img, image.Args{
			Type:    core1_0.ImageType2D,
			Width:   s.width,
			Height:  s.height,
			Depth:   1,
			Levels:  1,
			Format:  imageFormat,
			Usage:   imageUsage,
			Sharing: imageSharing,
		})
	})

	// create synchronization semaphores
	s.nextContext = 0
	s.contexts = make([]*Context, s.frames)
	for i := 0; i < s.frames; i++ {
		s.contexts[i] = NewContext(s.device, i)
	}
}

func (s *swapchain) Aquire(worker command.Worker) (*Context, error) {
	available := make(chan *Context)
	worker.Invoke(func() {
		defer close(available)

		// recreate if resized
		if s.resized {
			s.recreate()
			return
		}

		// get the next available context & update ring buffer index
		ctx := s.contexts[s.nextContext]
		s.nextContext = (s.nextContext + 1) % s.frames

		// aquiring the next frame
		idx, r, err := s.ptr.AcquireNextImage(1e9, ctx.ImageAvailable.Ptr(), nil)
		if err != nil {
			panic(err)
		}
		if r == khr_swapchain.VKErrorOutOfDate {
			s.recreate()
			return
		}

		// store frame index & return the context
		ctx.Index = idx
		available <- ctx
	})

	// this will actually block both the worker and the render thread until:
	//   - all pending work has been submitted
	//   - the image is aquired
	// is there any work for the render loop to do until the next frame is available?
	ctx := <-available
	if ctx == nil {
		return nil, fmt.Errorf("swapchain out of date")
	}
	return ctx, nil
}

func (s *swapchain) Present(worker command.Worker, ctx *Context) {
	if ctx.RenderComplete == nil {
		panic("context has no RenderComplete semaphore")
	}
	worker.Invoke(func() {
		// ideally there would be a better way to access the correct queue from the worker
		// however, this is the only place outside of the worker where we need to access the queue
		queue := s.device.Queue()
		s.ext.QueuePresent(queue.Ptr(), khr_swapchain.PresentInfo{
			WaitSemaphores: []core1_0.Semaphore{ctx.RenderComplete.Ptr()},
			Swapchains:     []khr_swapchain.Swapchain{s.ptr},
			ImageIndices:   []int{ctx.Index},
		})
	})
}

func (s *swapchain) Destroy() {
	for _, ctx := range s.contexts {
		ctx.Destroy()
	}
	s.contexts = nil

	if s.ptr != nil {
		s.ptr.Destroy(nil)
		s.ptr = nil
	}
}
