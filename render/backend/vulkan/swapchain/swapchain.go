package swapchain

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"
	"github.com/johanhenriksson/goworld/util"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Swapchain]

	Aquire() int
	Submit([]vk.CommandBuffer)
	Present()
	CurrentImage() image.T
	Resize(int, int)

	Count() int
	Image(int) image.T
}

type swapchain struct {
	ptr               vk.Swapchain
	device            device.T
	queue             vk.Queue
	surface           vk.Surface
	surfaceFormat     vk.SurfaceFormat
	images            []image.T
	fenceSwap         []sync.Fence
	semImageAvailable sync.Semaphore
	semRenderComplete sync.Semaphore
	currentImage      int
	swapCount         int
}

func New(device device.T, width, height, count int, surface vk.Surface, surfaceFormat vk.SurfaceFormat) T {
	// todo: surface format logic
	queue := device.GetQueue(0, vk.QueueFlags(vk.QueueGraphicsBit))

	s := &swapchain{
		device:        device,
		queue:         queue,
		surface:       surface,
		surfaceFormat: surfaceFormat,
		swapCount:     count,

		semImageAvailable: sync.NewSemaphore(device),
		semRenderComplete: sync.NewSemaphore(device),
	}

	// size according to framebuffer
	s.Resize(width, height)

	return s
}

func (s *swapchain) Ptr() vk.Swapchain {
	return s.ptr
}

func (s *swapchain) Resize(width, height int) {
	s.device.WaitIdle()

	if s.ptr != nil {
		vk.DestroySwapchain(s.device.Ptr(), s.ptr, nil)
	}

	swapInfo := vk.SwapchainCreateInfo{
		SType:           vk.StructureTypeSwapchainCreateInfo,
		Surface:         s.surface,
		MinImageCount:   uint32(s.swapCount),
		ImageFormat:     s.surfaceFormat.Format,
		ImageColorSpace: s.surfaceFormat.ColorSpace,
		ImageExtent: vk.Extent2D{
			Width:  uint32(width),
			Height: uint32(height),
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

	// get swapchain images
	swapImageCount := uint32(0)
	vk.GetSwapchainImages(s.device.Ptr(), s.ptr, &swapImageCount, nil)
	if swapImageCount != uint32(s.swapCount) {
		panic("failed to get the requested number of swapchain images")
	}

	images := make([]vk.Image, swapImageCount)
	vk.GetSwapchainImages(s.device.Ptr(), s.ptr, &swapImageCount, images)
	s.images = util.Map(images, func(i int, ptr vk.Image) image.T { return image.Wrap(s.device, ptr) })

	// set up fences for each backbuffer
	if len(s.fenceSwap) != int(swapImageCount) {
		for _, existing := range s.fenceSwap {
			existing.Destroy()
		}
		s.fenceSwap = make([]sync.Fence, swapImageCount)
		for i := range s.fenceSwap {
			s.fenceSwap[i] = sync.NewFence(s.device, true)
		}
	}
}

func (s *swapchain) Aquire() int {
	idx := uint32(0)
	vk.AcquireNextImage(s.device.Ptr(), s.ptr, vk.MaxUint64, s.semImageAvailable.Ptr(), nil, &idx)
	s.currentImage = int(idx)
	return s.currentImage
}

func (s *swapchain) Submit(commandBuffers []vk.CommandBuffer) {
	s.fenceSwap[s.currentImage].Wait()
	s.fenceSwap[s.currentImage].Reset()

	submitInfo := vk.SubmitInfo{
		SType:                vk.StructureTypeSubmitInfo,
		CommandBufferCount:   1,
		PCommandBuffers:      commandBuffers,
		WaitSemaphoreCount:   1,
		PWaitSemaphores:      []vk.Semaphore{s.semImageAvailable.Ptr()},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vk.Semaphore{s.semRenderComplete.Ptr()},
		PWaitDstStageMask: []vk.PipelineStageFlags{
			vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
		},
	}
	vk.QueueSubmit(s.queue, 1, []vk.SubmitInfo{submitInfo}, s.fenceSwap[s.currentImage].Ptr())
}

func (s *swapchain) Present() {
	presentInfo := vk.PresentInfo{
		SType:              vk.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vk.Semaphore{s.semRenderComplete.Ptr()},
		SwapchainCount:     1,
		PSwapchains:        []vk.Swapchain{s.ptr},
		PImageIndices:      []uint32{uint32(s.currentImage)},
	}
	vk.QueuePresent(s.queue, &presentInfo)
}

func (s *swapchain) CurrentImage() image.T {
	return s.images[s.currentImage]
}

func (s *swapchain) Count() int {
	return s.swapCount
}

func (s *swapchain) Image(idx int) image.T {
	return s.images[idx]
}

func (s *swapchain) Destroy() {
	for _, fence := range s.fenceSwap {
		fence.Destroy()
	}
	s.fenceSwap = nil

	s.semImageAvailable.Destroy()
	s.semImageAvailable = nil

	s.semRenderComplete.Destroy()
	s.semRenderComplete = nil

	if s.ptr != nil {
		vk.DestroySwapchain(s.device.Ptr(), s.ptr, nil)
		s.ptr = nil
	}
}
