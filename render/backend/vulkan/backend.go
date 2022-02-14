package vulkan

import (
	"fmt"

	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/command"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/device"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/framebuffer"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/image"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/instance"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/pipeline"
	"github.com/johanhenriksson/goworld/render/backend/vulkan/swapchain"

	"github.com/go-gl/glfw/v3.3/glfw"
	vk "github.com/vulkan-go/vulkan"
)

type VkVertex struct {
	X, Y, Z float32
	R, G, B float32
}

type T interface {
	Instance() instance.T
	Device() device.T
	Surface() vk.Surface
	Swapchain() swapchain.T
	Destroy()
	OutputPass() pipeline.Pass

	GlfwHints(window.Args) []window.GlfwHint
	GlfwSetup(*glfw.Window, window.Args) error

	Resize(int, int)
	CmdPool() command.Pool
	Framebuffer(int) framebuffer.T
	Aquire()
	Present()
	Submit([]vk.CommandBuffer)
}

type backend struct {
	appName   string
	deviceIdx int
	swapcount int
	instance  instance.T
	device    device.T
	surface   vk.Surface
	swapchain swapchain.T
	depth     image.T
	cmdpools  []command.Pool
	cmdpool   command.Pool
	swapviews []image.View
	buffers   []framebuffer.T
	buffer    framebuffer.T
	output    pipeline.Pass
	depthview image.View
}

func New(appName string, deviceIndex int) T {
	return &backend{
		appName:   appName,
		deviceIdx: deviceIndex,
		swapcount: 2,
	}
}

func (b *backend) Instance() instance.T   { return b.instance }
func (b *backend) Device() device.T       { return b.device }
func (b *backend) Surface() vk.Surface    { return b.surface }
func (b *backend) CmdPool() command.Pool  { return b.cmdpool }
func (b *backend) Swapchain() swapchain.T { return b.swapchain }

func (b *backend) GlfwHints(args window.Args) []window.GlfwHint {
	return []window.GlfwHint{
		{Hint: glfw.ClientAPI, Value: glfw.NoAPI},
	}
}

func (b *backend) GlfwSetup(w *glfw.Window, args window.Args) error {
	// initialize vulkan
	vk.SetGetInstanceProcAddr(glfw.GetVulkanGetInstanceProcAddress())
	if err := vk.Init(); err != nil {
		panic(err)
	}

	fmt.Println("window required extensions:", w.GetRequiredInstanceExtensions())

	// create instance * device
	b.instance = instance.New(b.appName)
	b.device = b.instance.GetDevice(b.deviceIdx)

	// surface
	surfPtr, err := w.CreateWindowSurface(b.instance.Ptr(), nil)
	if err != nil {
		panic(err)
	}

	b.surface = vk.SurfaceFromPointer(surfPtr)
	surfaceFormat := b.device.GetSurfaceFormats(b.surface)[0]

	// allocate swapchain
	width, height := w.GetFramebufferSize()
	b.swapchain = swapchain.New(b.device, width, height, b.swapcount, b.surface, surfaceFormat)

	// allocate a command pool for each swap image
	b.cmdpools = make([]command.Pool, b.swapcount)
	for i := range b.cmdpools {
		b.cmdpools[i] = command.NewPool(
			b.device,
			vk.CommandPoolCreateFlags(vk.CommandPoolCreateResetCommandBufferBit),
			vk.QueueFlags(vk.QueueGraphicsBit))
	}
	b.cmdpool = b.cmdpools[0]

	//
	// mega refactor below this point
	//

	// create output render pass
	b.output = pipeline.NewPass(b.device, &vk.RenderPassCreateInfo{
		SType:           vk.StructureTypeRenderPassCreateInfo,
		AttachmentCount: 2,
		PAttachments: []vk.AttachmentDescription{
			{
				Format:         surfaceFormat.Format,
				Samples:        vk.SampleCount1Bit,
				LoadOp:         vk.AttachmentLoadOpClear,
				StoreOp:        vk.AttachmentStoreOpStore,
				StencilLoadOp:  vk.AttachmentLoadOpDontCare,
				StencilStoreOp: vk.AttachmentStoreOpDontCare,
				InitialLayout:  vk.ImageLayoutUndefined,
				FinalLayout:    vk.ImageLayoutPresentSrc,
			},
			{
				Format:         b.device.GetDepthFormat(),
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

	// allocate a depth buffer
	depthFormat := b.device.GetDepthFormat()
	usage := vk.ImageUsageFlags(vk.ImageUsageDepthStencilAttachmentBit | vk.ImageUsageTransferSrcBit)
	b.depth = image.New2D(b.device, width, height, depthFormat, usage)
	b.depthview = b.depth.View(depthFormat, vk.ImageAspectFlags(vk.ImageAspectDepthBit|vk.ImageAspectStencilBit))

	// allocate a frame buffer for each swap image
	b.buffers = make([]framebuffer.T, b.swapcount)
	b.swapviews = make([]image.View, b.swapcount)
	for i := range b.buffers {
		colorview := b.swapchain.Image(i).View(surfaceFormat.Format, vk.ImageAspectFlags(vk.ImageAspectColorBit))
		b.swapviews[i] = colorview
		b.buffers[i] = framebuffer.New(
			b.device,
			width, height,
			b.output.Ptr(),
			[]image.View{
				colorview,
				b.depthview,
			},
		)
	}

	return nil
}

func (b *backend) Destroy() {
	b.output.Destroy()

	for _, buf := range b.buffers {
		buf.Destroy()
	}
	for _, view := range b.swapviews {
		view.Destroy()
	}
	b.depth.Destroy()
	b.depthview.Destroy()

	for _, pool := range b.cmdpools {
		pool.Destroy()
	}
	b.cmdpools = nil

	if b.swapchain != nil {
		b.swapchain.Destroy()
		b.swapchain = nil
	}
	if b.surface != nil {
		vk.DestroySurface(b.instance.Ptr(), b.surface, nil)
		b.surface = nil
	}
	if b.device != nil {
		b.device.Destroy()
		b.device = nil
	}
	if b.instance != nil {
		b.instance.Destroy()
		b.instance = nil
	}
}

func (b *backend) Framebuffer(idx int) framebuffer.T {
	return b.buffers[idx]
}

func (b *backend) OutputPass() pipeline.Pass {
	return b.output
}

func (b *backend) Resize(width, height int) {
	b.swapchain.Resize(width, height)
}

func (b *backend) Aquire() {
	index := b.swapchain.Aquire()
	b.cmdpool = b.cmdpools[index]
	b.buffer = b.buffers[index]
}

func (b *backend) Present() {
	b.swapchain.Present()
}

func (b *backend) Submit(cmdBuffers []vk.CommandBuffer) {
	b.swapchain.Submit(cmdBuffers)
}
