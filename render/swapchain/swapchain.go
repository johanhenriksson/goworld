package swapchain

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/render/sync"
	"github.com/johanhenriksson/goworld/util"

	"github.com/vkngwrapper/core/v2/common"
	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/khr_surface"
	"github.com/vkngwrapper/extensions/v2/khr_swapchain"
)

type T interface {
	device.Resource[khr_swapchain.Swapchain]

	Aquire() (*Context, error)
	Present(*Context)
	Resize(int, int)

	Images() []image.T
	SurfaceFormat() core1_0.Format
}

type swapchain struct {
	device     device.T
	ptr        khr_swapchain.Swapchain
	ext        khr_swapchain.Extension
	surface    khr_surface.Surface
	surfaceFmt khr_surface.SurfaceFormat
	images     []image.T
	frames     int
	width      int
	height     int
	resized    bool

	nextSemaphore  int
	imageAvailable []sync.Semaphore
	renderComplete []sync.Semaphore
}

func New(device device.T, frames, width, height int, surface khr_surface.Surface, surfaceFormat khr_surface.SurfaceFormat) T {
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

func (s *swapchain) Images() []image.T             { return s.images }
func (s *swapchain) SurfaceFormat() core1_0.Format { return core1_0.Format(s.surfaceFmt.Format) }

func (s *swapchain) Resize(width, height int) {
	// resizing actually happens the next time a frame is aquired
	s.width = width
	s.height = height
	s.resized = true
}

func (s *swapchain) recreate() {
	log.Println("recreating swapchain")

	// wait for device idle
	s.device.WaitIdle()

	// recreate swapchain resources
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
	s.images = util.Map(swapimages, func(img core1_0.Image) image.T {
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
	s.nextSemaphore = 0
	s.imageAvailable = make([]sync.Semaphore, s.frames)
	s.renderComplete = make([]sync.Semaphore, s.frames)
	for i := 0; i < s.frames; i++ {
		s.imageAvailable[i] = sync.NewSemaphore(s.device, fmt.Sprintf("ImageAvailable:%d", i))
		s.renderComplete[i] = sync.NewSemaphore(s.device, fmt.Sprintf("RenderComplete:%d", i))
	}
}

func (s *swapchain) Aquire() (*Context, error) {
	if s.resized {
		s.recreate()
		s.resized = false
		return nil, fmt.Errorf("swapchain out of date")
	}

	// create a context from the next set of available semaphores
	ctx := &Context{
		ImageAvailable: s.imageAvailable[s.nextSemaphore],
		RenderComplete: s.renderComplete[s.nextSemaphore],
	}

	idx, r, err := s.ptr.AcquireNextImage(common.NoTimeout, ctx.ImageAvailable.Ptr(), nil)
	if err != nil {
		panic(err)
	}
	if r == khr_swapchain.VKErrorOutOfDate {
		s.recreate()
		return nil, fmt.Errorf("swapchain out of date")
	}

	// store frame index
	ctx.Index = idx

	// update semaphore ring buffer index
	s.nextSemaphore = (s.nextSemaphore + 1) % s.frames

	return ctx, nil
}

func (s *swapchain) Present(ctx *Context) {
	if ctx.RenderComplete == nil {
		panic("context has no RenderComplete semaphore")
	}
	queue := s.device.Queue()
	s.ext.QueuePresent(queue.Ptr(), khr_swapchain.PresentInfo{
		WaitSemaphores: []core1_0.Semaphore{ctx.RenderComplete.Ptr()},
		Swapchains:     []khr_swapchain.Swapchain{s.Ptr()},
		ImageIndices:   []int{ctx.Index},
	})
}

func (s *swapchain) Destroy() {
	for i := range s.imageAvailable {
		s.imageAvailable[i].Destroy()
	}
	s.imageAvailable = nil
	for i := range s.renderComplete {
		s.renderComplete[i].Destroy()
	}
	s.renderComplete = nil

	if s.ptr != nil {
		s.ptr.Destroy(nil)
		s.ptr = nil
	}
}
