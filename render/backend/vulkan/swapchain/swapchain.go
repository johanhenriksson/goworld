package swapchain

import (
	"fmt"
	"log"

	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Swapchain]

	Aquire() (Context, error)
	Present()
	Resize(int, int)

	Images() []image.T
	SurfaceFormat() vk.Format
}

type swapchain struct {
	ptr        vk.Swapchain
	device     device.T
	queue      vk.Queue
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
	// todo: surface format logic
	queue := device.GetQueue(0, vk.QueueFlags(vk.QueueGraphicsBit))

	s := &swapchain{
		device:     device,
		queue:      queue,
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
	log.Println("Resize swapchain to", width, "x", height)
}

func (s *swapchain) recreate() {
	log.Println("Recreating swapchain")
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
		ImageUsage:       vk.ImageUsageFlags(vk.ImageUsageColorAttachmentBit),
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
	s.images = util.Map(swapimages, func(img vk.Image) image.T { return image.Wrap(s.device, img) })

	for i := range s.contexts {
		// destroy existing
		s.contexts[i].Destroy()

		s.contexts[i] = Context{
			Index:          i,
			ImageAvailable: sync.NewSemaphore(s.device),
			RenderComplete: sync.NewSemaphore(s.device),
		}
	}

	// this ensures the first call to Aquire works properly
	s.current = -1
}

func (s *swapchain) Aquire() (Context, error) {
	idx := uint32(0)
	next := s.contexts[(s.current+1)%s.frames]
	r := vk.AcquireNextImage(s.device.Ptr(), s.ptr, 1e9, next.ImageAvailable.Ptr(), nil, &idx)
	if r == vk.ErrorOutOfDate {
		s.recreate()
		return Context{}, fmt.Errorf("swapchain out of date")
	}
	s.current = int(idx)
	return s.contexts[s.current], nil
}

// Present executes the final output render pass and presents the image to the screen.
// When calling, the final render pass command buffer should be submitted to worker 0 of the current context.
func (s *swapchain) Present() {
	ctx := s.contexts[s.current]
	presentInfo := vk.PresentInfo{
		SType:              vk.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vk.Semaphore{ctx.RenderComplete.Ptr()},
		SwapchainCount:     1,
		PSwapchains:        []vk.Swapchain{s.ptr},
		PImageIndices:      []uint32{uint32(s.current)},
	}

	vk.QueuePresent(s.queue, &presentInfo)
	if s.resized {
		s.recreate()
		s.resized = false
	}
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
