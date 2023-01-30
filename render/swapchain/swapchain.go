package swapchain

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/image"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Swapchain]

	Aquire() (Context, error)
	Resize(int, int)

	Images() []image.T
	SurfaceFormat() vk.Format
}

type swapchain struct {
	ptr        vk.Swapchain
	device     device.T
	surface    vk.Surface
	surfaceFmt vk.SurfaceFormat
	images     []image.T
	current    int
	frames     int
	width      int
	height     int
	resized    bool

	contexts []Context
}

func New(device device.T, frames, width, height int, surface vk.Surface, surfaceFormat vk.SurfaceFormat) T {
	s := &swapchain{
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

func (s *swapchain) Ptr() vk.Swapchain {
	return s.ptr
}

func (s *swapchain) Images() []image.T        { return s.images }
func (s *swapchain) SurfaceFormat() vk.Format { return s.surfaceFmt.Format }

func (s *swapchain) Resize(width, height int) {
	s.width = width
	s.height = height
	s.resized = true
}

func (s *swapchain) recreate() {
	log.Println("recreating swapchain")
	s.device.WaitIdle()

	if s.ptr != nil {
		vk.DestroySwapchain(s.device.Ptr(), s.ptr, nil)
	}

	swapInfo := vk.SwapchainCreateInfo{
		SType:           vk.StructureTypeSwapchainCreateInfo,
		Surface:         s.surface,
		MinImageCount:   uint32(s.frames),
		ImageFormat:     s.surfaceFmt.Format,
		ImageColorSpace: s.surfaceFmt.ColorSpace,
		ImageExtent: vk.Extent2D{
			Width:  uint32(s.width),
			Height: uint32(s.height),
		},
		ImageArrayLayers: 1,
		ImageUsage:       vk.ImageUsageFlags(vk.ImageUsageColorAttachmentBit | vk.ImageUsageTransferSrcBit),
		ImageSharingMode: vk.SharingModeExclusive,
		PresentMode:      vk.PresentModeFifo,
		PreTransform:     vk.SurfaceTransformIdentityBit,
		CompositeAlpha:   vk.CompositeAlphaOpaqueBit,
		Clipped:          vk.True,
	}

	var chain vk.Swapchain
	r := vk.CreateSwapchain(s.device.Ptr(), &swapInfo, nil, &chain)
	if r != vk.Success {
		panic("failed to create swapchain")
	}
	s.ptr = chain

	swapImageCount := uint32(0)
	vk.GetSwapchainImages(s.device.Ptr(), s.ptr, &swapImageCount, nil)
	if swapImageCount != uint32(s.frames) {
		panic("failed to get the requested number of swapchain images")
	}

	swapimages := make([]vk.Image, swapImageCount)
	vk.GetSwapchainImages(s.device.Ptr(), s.ptr, &swapImageCount, swapimages)
	s.images = util.Map(swapimages, func(img vk.Image) image.T {
		return image.Wrap(s.device, img, image.Args{
			Width:  s.width,
			Height: s.height,
			Depth:  1,
			Format: s.surfaceFmt.Format,
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

	idx := uint32(0)
	timeoutNs := uint64(1e9)
	nextFrame := s.contexts[(s.current+1)%s.frames]

	r := vk.AcquireNextImage(s.device.Ptr(), s.ptr, timeoutNs, nextFrame.ImageAvailable.Ptr(), nil, &idx)
	if r == vk.ErrorOutOfDate {
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
		vk.DestroySwapchain(s.device.Ptr(), s.ptr, nil)
		s.ptr = nil
	}
}
