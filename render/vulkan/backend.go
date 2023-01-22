package vulkan

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/swapchain"
	"github.com/johanhenriksson/goworld/render/vulkan/instance"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	Instance() instance.T
	Device() device.T
	Swapchain() swapchain.T
	Frames() int
	Width() int
	Height() int
	Destroy()
	Window(WindowArgs) (Window, error)

	Aquire() (swapchain.Context, error)
	Present()

	Worker(int) command.Worker
	Transferer() command.Worker

	Meshes() cache.MeshCache
	Textures() cache.TextureCache
}

type backend struct {
	appName   string
	frames    int
	width     int
	height    int
	instance  instance.T
	device    device.T
	swapchain swapchain.T

	transfer command.Worker
	workers  []command.Worker
	windows  []Window

	meshes   cache.MeshCache
	textures cache.TextureCache
}

func New(appName string, deviceIndex int) T {
	// create instance * device
	frames := 3
	instance := instance.New(appName)
	device := instance.GetDevice(deviceIndex)

	// transfer worker
	transfer := command.NewWorker(device, vk.QueueFlags(vk.QueueTransferBit))

	// per frame graphics workers
	workers := make([]command.Worker, frames)
	for i := 0; i < frames; i++ {
		workers[i] = command.NewWorker(device, vk.QueueFlags(vk.QueueGraphicsBit))
	}

	// init global descriptor pool
	descriptor.InitGlobalPool(device)

	// init caches
	meshes := cache.NewMeshCache(device, transfer)
	textures := cache.NewTextureCache(device, transfer)

	return &backend{
		appName: appName,
		frames:  frames,
		windows: []Window{},

		device:   device,
		instance: instance,
		transfer: transfer,
		workers:  workers,
		meshes:   meshes,
		textures: textures,
	}
}

func (b *backend) Instance() instance.T   { return b.instance }
func (b *backend) Device() device.T       { return b.device }
func (b *backend) Swapchain() swapchain.T { return b.swapchain }
func (b *backend) Frames() int            { return b.frames }
func (b *backend) Width() int             { return b.width }
func (b *backend) Height() int            { return b.height }

func (b *backend) Meshes() cache.MeshCache      { return b.meshes }
func (b *backend) Textures() cache.TextureCache { return b.textures }

func (b *backend) Transferer() command.Worker {
	return b.transfer
}

func (b *backend) Worker(frame int) command.Worker {
	return b.workers[frame]
}

func (b *backend) Window(args WindowArgs) (Window, error) {
	w, err := NewWindow(b, args)
	if err != nil {
		return nil, fmt.Errorf("failed to create window: %w", err)
	}

	b.swapchain = w.Swapchain()
	b.width, b.height = w.Size()
	return w, nil
}

func (b *backend) Destroy() {
	// clean up caches
	b.meshes.Destroy()
	b.textures.Destroy()

	// clean up global descriptor pool
	descriptor.DestroyGlobalPool()

	for i := 0; i < b.frames; i++ {
		b.workers[i].Destroy()
	}
	b.transfer.Destroy()

	for _, wnd := range b.windows {
		wnd.Destroy()
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

func (b *backend) Resize(width, height int) {
	b.width = width
	b.height = height
	b.swapchain.Resize(width, height)
}

func (b *backend) Aquire() (swapchain.Context, error) {
	return b.swapchain.Aquire()
}

func (b *backend) Present() {
	b.swapchain.Present()
	b.meshes.Tick()
	b.textures.Tick()
}
