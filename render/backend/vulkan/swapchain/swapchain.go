package swapchain

import (
	"fmt"
	"log"

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

	Aquire() (Context, error)
	Present()
	Resize(int, int)

	Output() pipeline.Pass
	Context() Context
	Count() int
}

type swapchain struct {
	ptr        vk.Swapchain
	device     device.T
	queue      vk.Queue
	surface    vk.Surface
	surfaceFmt vk.SurfaceFormat
	output     pipeline.Pass
	current    int
	buffers    int
	width      int
	height     int
	resized    bool

	contexts []Context
}

func New(device device.T, buffers, width, height int, surface vk.Surface, surfaceFormat vk.SurfaceFormat) T {
	// todo: surface format logic
	queue := device.GetQueue(0, vk.QueueFlags(vk.QueueGraphicsBit))

	s := &swapchain{
		device:     device,
		queue:      queue,
		surface:    surface,
		surfaceFmt: surfaceFormat,
		buffers:    buffers,
		contexts:   make([]Context, buffers),
		width:      width,
		height:     height,
	}

	s.output = s.createOutputPass()

	s.recreate()

	return s
}

func (s *swapchain) Ptr() vk.Swapchain {
	return s.ptr
}

func (s *swapchain) Output() pipeline.Pass {
	return s.output
}

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

	// query max surface size
	// caps := s.device.GetSurfaceCapabilities(s.surface)
	// caps.MaxImageExtent.Deref()
	// s.width = int(caps.MaxImageExtent.Width)
	// s.height = int(caps.MaxImageExtent.Height)

	swapInfo := vk.SwapchainCreateInfo{
		SType:           vk.StructureTypeSwapchainCreateInfo,
		Surface:         s.surface,
		MinImageCount:   uint32(s.buffers),
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
	if swapImageCount != uint32(s.buffers) {
		panic("failed to get the requested number of swapchain images")
	}
	s.buffers = int(swapImageCount)

	images := make([]vk.Image, swapImageCount)
	vk.GetSwapchainImages(s.device.Ptr(), s.ptr, &swapImageCount, images)

	// depth buffer
	// todo: destroy existing
	depthFormat := s.device.GetDepthFormat()
	depthUsage := vk.ImageUsageFlags(vk.ImageUsageDepthStencilAttachmentBit | vk.ImageUsageTransferSrcBit)

	for i := range s.contexts {
		// destroy existing
		s.contexts[i].Destroy()

		color := image.Wrap(s.device, images[i])
		colorview := color.View(s.surfaceFmt.Format, vk.ImageAspectFlags(vk.ImageAspectColorBit))

		depth := image.New2D(s.device, s.width, s.height, depthFormat, depthUsage)
		depthview := depth.View(depthFormat, vk.ImageAspectFlags(vk.ImageAspectDepthBit|vk.ImageAspectStencilBit))

		s.contexts[i] = Context{
			Index:     i,
			Color:     color,
			ColorView: colorview,
			Depth:     depth,
			DepthView: depthview,
			Framebuffer: framebuffer.New(
				s.device,
				s.width, s.height,
				s.output.Ptr(),
				[]image.View{
					colorview,
					depthview,
				},
			),
			Workers: command.Workers{
				command.NewWorker(s.device, vk.QueueFlags(vk.QueueGraphicsBit)),
			},
			Output:         s.output,
			Width:          s.width,
			Height:         s.height,
			ImageAvailable: sync.NewSemaphore(s.device),
			RenderComplete: sync.NewSemaphore(s.device),
		}
	}

	// this ensures the first call to Aquire works properly
	s.current = -1
}

func (s *swapchain) Aquire() (Context, error) {
	idx := uint32(0)
	next := s.contexts[(s.current+1)%s.buffers]
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
	ctx := s.Context()
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

func (s *swapchain) Context() Context {
	return s.contexts[s.current]
}

func (s *swapchain) Count() int {
	return s.buffers
}

func (s *swapchain) Destroy() {
	for _, context := range s.contexts {
		context.Destroy()
	}
	s.contexts = nil

	s.output.Destroy()

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
				StoreOp:        vk.AttachmentStoreOpStore,
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
		PDependencies: []vk.SubpassDependency{
			{
				SrcSubpass:      0,
				DstSubpass:      0,
				SrcStageMask:    vk.PipelineStageFlags(vk.PipelineStageEarlyFragmentTestsBit | vk.PipelineStageLateFragmentTestsBit),
				SrcAccessMask:   0,
				DstStageMask:    vk.PipelineStageFlags(vk.PipelineStageEarlyFragmentTestsBit | vk.PipelineStageLateFragmentTestsBit),
				DstAccessMask:   vk.AccessFlags(vk.AccessDepthStencilAttachmentWriteBit),
				DependencyFlags: vk.DependencyFlags(vk.DependencyByRegionBit),
			},
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
