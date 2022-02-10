package swapchain

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/fence"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/semaphore"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource
	Ptr() vk.Swapchain
	Resize(int, int)

	Aquire()
	Submit([]vk.CommandBuffer)
	Present()
	NextImage() vk.Image
}

type swapchain struct {
	ptr               vk.Swapchain
	device            device.T
	queue             vk.Queue
	surface           vk.Surface
	surfaceFormat     vk.SurfaceFormat
	images            []vk.Image
	fenceSwap         fence.T
	semImageAvailable semaphore.T
	semRenderComplete semaphore.T
	nextImage         uint32
}

func New(window *glfw.Window, device device.T, surface vk.Surface) T {
	// todo: surface format logic
	surfaceFormat := device.GetSurfaceFormats(surface)[0]

	s := &swapchain{
		device:        device,
		queue:         device.GetQueue(0, 0),
		surface:       surface,
		surfaceFormat: surfaceFormat,

		semImageAvailable: semaphore.New(device),
		semRenderComplete: semaphore.New(device),
		fenceSwap:         fence.New(device),
	}

	// size according to framebuffer
	width, height := window.GetFramebufferSize()
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
		MinImageCount:   2,
		ImageFormat:     s.surfaceFormat.Format,
		ImageColorSpace: s.surfaceFormat.ColorSpace,
		ImageExtent: vk.Extent2D{
			Width:  uint32(width),
			Height: uint32(height),
		},
		ImageArrayLayers: 1,
		ImageUsage:       vk.ImageUsageFlags(vk.ImageUsageColorAttachmentBit),
		PresentMode:      vk.PresentModeFifo,
		PreTransform:     vk.SurfaceTransformIdentityBit,
		CompositeAlpha:   vk.CompositeAlphaOpaqueBit,
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
	images := make([]vk.Image, swapImageCount)
	vk.GetSwapchainImages(s.device.Ptr(), s.ptr, &swapImageCount, images)
	s.images = images
}

func (s *swapchain) Aquire() {
	idx := uint32(0)
	vk.AcquireNextImage(s.device.Ptr(), s.ptr, vk.MaxUint64, s.semImageAvailable.Ptr(), nil, &idx)
	s.nextImage = idx
}

func (s *swapchain) Submit(commandBuffers []vk.CommandBuffer) {
	submitInfo := vk.SubmitInfo{
		SType:                vk.StructureTypeSubmitInfo,
		CommandBufferCount:   0,
		WaitSemaphoreCount:   1,
		PWaitSemaphores:      []vk.Semaphore{s.semImageAvailable.Ptr()},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vk.Semaphore{s.semRenderComplete.Ptr()},
		PWaitDstStageMask: []vk.PipelineStageFlags{
			vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
		},
	}
	vk.QueueSubmit(s.queue, 1, []vk.SubmitInfo{submitInfo}, s.fenceSwap.Ptr())
}

func (s *swapchain) Present() {
	vk.WaitForFences(s.device.Ptr(), 1, []vk.Fence{s.fenceSwap.Ptr()}, vk.Bool32(0), vk.MaxUint64)
	vk.ResetFences(s.device.Ptr(), 1, []vk.Fence{s.fenceSwap.Ptr()})

	presentInfo := vk.PresentInfo{
		SType:              vk.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vk.Semaphore{s.semRenderComplete.Ptr()},
		SwapchainCount:     1,
		PSwapchains:        []vk.Swapchain{s.ptr},
		PImageIndices:      []uint32{s.nextImage},
	}
	vk.QueuePresent(s.queue, &presentInfo)
}

func (s *swapchain) NextImage() vk.Image {
	return s.images[s.nextImage]
}

func (s *swapchain) Destroy() {
	s.fenceSwap.Destroy()
	s.semImageAvailable.Destroy()
	s.semRenderComplete.Destroy()

	if s.ptr != nil {
		vk.DestroySwapchain(s.device.Ptr(), s.ptr, nil)
		s.ptr = nil
	}
}
