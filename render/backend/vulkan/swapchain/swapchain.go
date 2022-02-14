package swapchain

import (
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/framebuffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/sync"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	device.Resource[vk.Swapchain]

	Aquire() Context
	Present()
	Resize(int, int)

	Output() pipeline.Pass
	Context() Context
	Count() int
}

type Context struct {
	Index       int
	Color       image.T
	Depth       image.T
	ColorView   image.View
	DepthView   image.View
	Framebuffer framebuffer.T
	Workers     command.Workers
}

type swapchain struct {
	ptr               vk.Swapchain
	device            device.T
	queue             vk.Queue
	surface           vk.Surface
	surfaceFormat     vk.SurfaceFormat
	depth             image.T
	depthview         image.View
	semImageAvailable sync.Semaphore
	semRenderComplete sync.Semaphore
	current           int
	swapCount         int
	output            pipeline.Pass

	contexts []Context
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
		contexts:      make([]Context, count),

		semImageAvailable: sync.NewSemaphore(device),
		semRenderComplete: sync.NewSemaphore(device),
	}

	s.output = s.createOutputPass()

	// size according to framebuffer
	s.Resize(width, height)

	return s
}

func (s *swapchain) Ptr() vk.Swapchain {
	return s.ptr
}

func (s *swapchain) Output() pipeline.Pass {
	return s.output
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

	swapImageCount := uint32(0)
	vk.GetSwapchainImages(s.device.Ptr(), s.ptr, &swapImageCount, nil)
	if swapImageCount != uint32(s.swapCount) {
		panic("failed to get the requested number of swapchain images")
	}

	images := make([]vk.Image, swapImageCount)
	vk.GetSwapchainImages(s.device.Ptr(), s.ptr, &swapImageCount, images)

	// depth buffer
	// todo: destroy existing
	depthFormat := s.device.GetDepthFormat()
	usage := vk.ImageUsageFlags(vk.ImageUsageDepthStencilAttachmentBit | vk.ImageUsageTransferSrcBit)
	s.depth = image.New2D(s.device, width, height, depthFormat, usage)
	s.depthview = s.depth.View(depthFormat, vk.ImageAspectFlags(vk.ImageAspectDepthBit|vk.ImageAspectStencilBit))

	for i := range s.contexts {
		// todo: destroy existing

		color := image.Wrap(s.device, images[i])
		colorview := color.View(s.surfaceFormat.Format, vk.ImageAspectFlags(vk.ImageAspectColorBit))

		s.contexts[i] = Context{
			Index:     i,
			Color:     color,
			ColorView: colorview,
			Depth:     s.depth,
			DepthView: s.depthview,
			Framebuffer: framebuffer.New(
				s.device,
				width, height,
				s.output.Ptr(),
				[]image.View{
					colorview,
					s.depthview,
				},
			),
			Workers: command.Workers{
				command.NewWorker(s.device),
			},
		}
	}
}

func (s *swapchain) Aquire() Context {
	idx := uint32(0)
	vk.AcquireNextImage(s.device.Ptr(), s.ptr, vk.MaxUint64, s.semImageAvailable.Ptr(), nil, &idx)
	s.current = int(idx)
	return s.contexts[s.current]
}

func (s *swapchain) Present() {
	s.Context().Workers[0].Submit(command.SubmitInfo{
		Queue:  s.queue,
		Wait:   []sync.Semaphore{s.semImageAvailable},
		Signal: []sync.Semaphore{s.semRenderComplete},
		WaitMask: []vk.PipelineStageFlags{
			vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
		},
	})

	presentInfo := vk.PresentInfo{
		SType:              vk.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vk.Semaphore{s.semRenderComplete.Ptr()},
		SwapchainCount:     1,
		PSwapchains:        []vk.Swapchain{s.ptr},
		PImageIndices:      []uint32{uint32(s.current)},
	}
	vk.QueuePresent(s.queue, &presentInfo)
}

func (s *swapchain) Context() Context {
	return s.contexts[s.current]
}

func (s *swapchain) Count() int {
	return s.swapCount
}

func (s *swapchain) Destroy() {
	for _, context := range s.contexts {
		context.Destroy()
	}
	s.contexts = nil

	s.depthview.Destroy()
	s.depth.Destroy()

	s.output.Destroy()

	s.semImageAvailable.Destroy()
	s.semImageAvailable = nil

	s.semRenderComplete.Destroy()
	s.semRenderComplete = nil

	if s.ptr != nil {
		vk.DestroySwapchain(s.device.Ptr(), s.ptr, nil)
		s.ptr = nil
	}
}

func (s *swapchain) createOutputPass() pipeline.Pass {
	return pipeline.NewPass(s.device, &vk.RenderPassCreateInfo{
		SType:           vk.StructureTypeRenderPassCreateInfo,
		AttachmentCount: 2,
		PAttachments: []vk.AttachmentDescription{
			{
				Format:         s.surfaceFormat.Format,
				Samples:        vk.SampleCount1Bit,
				LoadOp:         vk.AttachmentLoadOpClear,
				StoreOp:        vk.AttachmentStoreOpStore,
				StencilLoadOp:  vk.AttachmentLoadOpDontCare,
				StencilStoreOp: vk.AttachmentStoreOpDontCare,
				InitialLayout:  vk.ImageLayoutUndefined,
				FinalLayout:    vk.ImageLayoutPresentSrc,
			},
			{
				Format:         s.device.GetDepthFormat(),
				Samples:        vk.SampleCount1Bit,
				LoadOp:         vk.AttachmentLoadOpClear,
				StoreOp:        vk.AttachmentStoreOpDontCare,
				StencilLoadOp:  vk.AttachmentLoadOpDontCare,
				StencilStoreOp: vk.AttachmentStoreOpDontCare,
				InitialLayout:  vk.ImageLayoutUndefined,
				FinalLayout:    vk.ImageLayoutDepthStencilAttachmentOptimal,
			},
		},
		SubpassCount: 1,
		PSubpasses: []vk.SubpassDescription{
			{
				PipelineBindPoint:    vk.PipelineBindPointGraphics,
				InputAttachmentCount: 0,
				ColorAttachmentCount: 1,
				PColorAttachments: []vk.AttachmentReference{
					{
						Attachment: 0,
						Layout:     vk.ImageLayoutColorAttachmentOptimal,
					},
				},
				PDepthStencilAttachment: &vk.AttachmentReference{
					Attachment: 1,
					Layout:     vk.ImageLayoutDepthStencilAttachmentOptimal,
				},
			},
		},
		DependencyCount: 0,
		PDependencies:   []vk.SubpassDependency{
			// {
			// 	SrcSubpass:      0,
			// 	DstSubpass:      0,
			// 	SrcStageMask:    vk.PipelineStageFlags(vk.PipelineStageBottomOfPipeBit),
			// 	DstStageMask:    vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
			// 	SrcAccessMask:   vk.AccessFlags(vk.AccessMemoryReadBit),
			// 	DstAccessMask:   vk.AccessFlags(vk.AccessColorAttachmentReadBit | vk.AccessColorAttachmentWriteBit),
			// 	DependencyFlags: vk.DependencyFlags(vk.DependencyByRegionBit),
			// },
			// {
			// 	SrcSubpass:      0,
			// 	DstSubpass:      0,
			// 	SrcStageMask:    vk.PipelineStageFlags(vk.PipelineStageBottomOfPipeBit),
			// 	DstStageMask:    vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
			// 	SrcAccessMask:   vk.AccessFlags(vk.AccessMemoryReadBit),
			// 	DstAccessMask:   vk.AccessFlags(vk.AccessColorAttachmentReadBit | vk.AccessColorAttachmentWriteBit),
			// 	DependencyFlags: vk.DependencyFlags(vk.DependencyByRegionBit),
			// },
		},
	})
}

func (c Context) Destroy() {
	c.ColorView.Destroy()
	c.Framebuffer.Destroy()
	c.Workers[0].Destroy()
}
