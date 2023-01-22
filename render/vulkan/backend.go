package vulkan

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vulkan/instance"

	vk "github.com/vulkan-go/vulkan"
)

type T interface {
	Instance() instance.T
	Device() device.T
	Destroy()
	Window(WindowArgs) (Window, error)

	Worker(int) command.Worker
	Transferer() command.Worker

	Meshes() cache.MeshCache
	Textures() cache.TextureCache
}

type backend struct {
	appName  string
	frames   int
	instance instance.T
	device   device.T

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

func (b *backend) Instance() instance.T { return b.instance }
func (b *backend) Device() device.T     { return b.device }
func (b *backend) Frames() int          { return b.frames }

func (b *backend) Meshes() cache.MeshCache      { return b.meshes }
func (b *backend) Textures() cache.TextureCache { return b.textures }

func (b *backend) Transferer() command.Worker {
	return b.transfer
}

func (b *backend) Worker(frame int) command.Worker {
	return b.workers[frame]
}

func (b *backend) Window(args WindowArgs) (Window, error) {
	args.Frames = b.frames
	w, err := NewWindow(b, args)
	if err != nil {
		return nil, fmt.Errorf("failed to create window: %w", err)
	}
	b.windows = append(b.windows, w)
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
