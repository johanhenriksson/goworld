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
	Present(command.CommandFn)
	Resize(int, int)

	Output() pipeline.Pass
	Context() Context
	Count() int
}

type swapchain struct {
	ptr            vk.Swapchain
	device         device.T
	queue          vk.Queue
	surface        vk.Surface
	surfaceFmt     vk.SurfaceFormat
	depth          image.T
	depthview      image.View
	imageAvailable sync.Semaphore
	renderComplete sync.Semaphore
	current        int
	swapCount      int
	output         pipeline.Pass

	contexts []Context
}

func New(device device.T, width, height, count int, surface vk.Surface, surfaceFormat vk.SurfaceFormat) T {
	// todo: surface format logic
	queue := device.GetQueue(0, vk.QueueFlags(vk.QueueGraphicsBit))

	s := &swapchain{
		device:     device,
		queue:      queue,
		surface:    surface,
		surfaceFmt: surfaceFormat,
		swapCount:  count,
		contexts:   make([]Context, count),

		imageAvailable: sync.NewSemaphore(device),
		renderComplete: sync.NewSemaphore(device),
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
		ImageFormat:     s.surfaceFmt.Format,
		ImageColorSpace: s.surfaceFmt.ColorSpace,
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
	if s.depth != nil {
		s.depth.Destroy()
	}
	if s.depthview != nil {
		s.depthview.Destroy()
	}
	depthFormat := s.device.GetDepthFormat()
	usage := vk.ImageUsageFlags(vk.ImageUsageDepthStencilAttachmentBit | vk.ImageUsageTransferSrcBit)
	s.depth = image.New2D(s.device, width, height, depthFormat, usage)
	s.depthview = s.depth.View(depthFormat, vk.ImageAspectFlags(vk.ImageAspectDepthBit|vk.ImageAspectStencilBit))

	for i := range s.contexts {
		// destroy existing
		s.contexts[i].Destroy()

		color := image.Wrap(s.device, images[i])
		colorview := color.View(s.surfaceFmt.Format, vk.ImageAspectFlags(vk.ImageAspectColorBit))

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
			Output: s.output,
			Width:  width,
			Height: height,
		}
	}
}

func (s *swapchain) Aquire() Context {
	idx := uint32(0)
	vk.AcquireNextImage(s.device.Ptr(), s.ptr, vk.MaxUint64, s.imageAvailable.Ptr(), nil, &idx)
	s.current = int(idx)
	return s.contexts[s.current]
}

// Present executes the final output render pass and presents the image to the screen.
// When calling, the final render pass command buffer should be submitted to worker 0 of the current context.
func (s *swapchain) Present(cmdf command.CommandFn) {
	context := s.Context()
	worker := context.Workers[0]

	worker.Queue(func(cmd command.Buffer) {
		clearValues := make([]vk.ClearValue, 2)
		clearValues[1].SetDepthStencil(1, 0)
		clearValues[0].SetColor([]float32{
			0.2, 0.2, 0.2, 0.2,
		})

		vk.CmdBeginRenderPass(cmd.Ptr(), &vk.RenderPassBeginInfo{
			SType:       vk.StructureTypeRenderPassBeginInfo,
			RenderPass:  s.Output().Ptr(),
			Framebuffer: context.Framebuffer.Ptr(),
			RenderArea: vk.Rect2D{
				Offset: vk.Offset2D{},
				Extent: vk.Extent2D{
					Width:  uint32(context.Width),
					Height: uint32(context.Height),
				},
			},
			ClearValueCount: 2,
			PClearValues:    clearValues,
		}, vk.SubpassContentsInline)

		vk.CmdSetViewport(cmd.Ptr(), 0, 1, []vk.Viewport{
			{
				Width:  float32(context.Width),
				Height: float32(context.Height),
			},
		})
		vk.CmdSetScissor(cmd.Ptr(), 0, 1, []vk.Rect2D{
			{
				Offset: vk.Offset2D{},
				Extent: vk.Extent2D{
					Width:  uint32(context.Width),
					Height: uint32(context.Height),
				},
			},
		})

		// user draw calls
		cmdf(cmd)

		vk.CmdEndRenderPass(cmd.Ptr())
	})

	s.Context().Workers[0].Submit(command.SubmitInfo{
		Queue:  s.queue,
		Wait:   []sync.Semaphore{s.imageAvailable},
		Signal: []sync.Semaphore{s.renderComplete},
		WaitMask: []vk.PipelineStageFlags{
			vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit),
		},
	})

	presentInfo := vk.PresentInfo{
		SType:              vk.StructureTypePresentInfo,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vk.Semaphore{s.renderComplete.Ptr()},
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

	s.imageAvailable.Destroy()
	s.imageAvailable = nil

	s.renderComplete.Destroy()
	s.renderComplete = nil

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
				Format:         s.surfaceFmt.Format,
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
