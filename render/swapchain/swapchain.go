package swapchain

import (
	"fmt"
	"log"
	"time"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/util"

	"github.com/vkngwrapper/core/v2/core1_0"
	"github.com/vkngwrapper/extensions/v2/khr_surface"
	"github.com/vkngwrapper/extensions/v2/khr_swapchain"
)

type T interface {
	device.Resource[khr_swapchain.Swapchain]

	Aquire() (Context, error)
	Resize(int, int)

	Images() []image.T
	SurfaceFormat() core1_0.Format
}

type swapchain struct {
	ext        khr_swapchain.Extension
	ptr        khr_swapchain.Swapchain
	device     device.T
	surface    khr_surface.Surface
	surfaceFmt khr_surface.SurfaceFormat
	images     []image.T
	current    int
	frames     int
	width      int
	height     int
	resized    bool

	contexts []Context
}

func New(device device.T, frames, width, height int, surface khr_surface.Surface, surfaceFormat khr_surface.SurfaceFormat) T {
	s := &swapchain{
		ext:        khr_swapchain.CreateExtensionFromDevice(device.Ptr()),
		device:     device,
		surface:    surface,
		surfaceFmt: surfaceFormat,
		frames:     frames,
		contexts:   make([]Context, frames),
		width:      width,
		height:     height,
	}

	s.recreate()

	return s
}

func (s *swapchain) Ptr() khr_swapchain.Swapchain {
	return s.ptr
}

func (s *swapchain) Images() []image.T             { return s.images }
func (s *swapchain) SurfaceFormat() core1_0.Format { return core1_0.Format(s.surfaceFmt.Format) }

func (s *swapchain) Resize(width, height int) {
	s.width = width
	s.height = height
	s.resized = true
}

func (s *swapchain) recreate() {
	log.Println("recreating swapchain")
	s.device.WaitIdle()

	if s.ptr != nil {
		s.ptr.Destroy(nil)
	}

	swapInfo := khr_swapchain.SwapchainCreateInfo{
		Surface:         s.surface,
		MinImageCount:   s.frames,
		ImageFormat:     core1_0.Format(s.surfaceFmt.Format),
		ImageColorSpace: khr_surface.ColorSpace(s.surfaceFmt.ColorSpace),
		ImageExtent: core1_0.Extent2D{
			Width:  s.width,
			Height: s.height,
		},
		ImageArrayLayers: 1,
		ImageUsage:       core1_0.ImageUsageColorAttachment | core1_0.ImageUsageTransferSrc,
		ImageSharingMode: core1_0.SharingModeExclusive,
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

	swapimages, _, err := chain.SwapchainImages()
	if err != nil {
		panic(err)
	}
	if len(swapimages) != s.frames {
		panic("failed to get the requested number of swapchain images")
	}
	s.images = util.Map(swapimages, func(img core1_0.Image) image.T {
		return image.Wrap(s.device, img, image.Args{
			Width:  s.width,
			Height: s.height,
			Depth:  1,
			Format: core1_0.Format(s.surfaceFmt.Format),
		})
	})

	for i := range s.contexts {
		// destroy existing
		s.contexts[i].Destroy()
		s.contexts[i] = newContext(s.device, i)
	}

	// this ensures the first call to Aquire works properly
	s.current = -1
}

func (s *swapchain) Aquire() (Context, error) {
	if s.resized {
		log.Println("aquire triggered swapchain recreation")
		s.recreate()
		s.resized = false
		return Context{}, fmt.Errorf("swapchain out of date")
	}

	nextFrame := s.contexts[(s.current+1)%s.frames]
	idx, r, err := s.ptr.AcquireNextImage(time.Second, nextFrame.ImageAvailable.Ptr(), nil)
	if err != nil {
		panic(err)
	}
	if r == khr_swapchain.VKErrorOutOfDate {
		s.recreate()
		return Context{}, fmt.Errorf("swapchain out of date")
	}

	s.current = int(idx)
	return s.contexts[s.current], nil
}

func (s *swapchain) Destroy() {
	for _, context := range s.contexts {
		context.Destroy()
	}
	s.contexts = nil

	if s.ptr != nil {
		s.ptr.Destroy(nil)
		s.ptr = nil
	}
}
