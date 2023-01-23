package vulkan

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
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
	windows  []Window

	transfer command.Worker
	workers  []command.Worker

	meshes   cache.MeshCache
	textures cache.TextureCache
}

func New(appName string, deviceIndex int) T {
	// create instance * device
	frames := 2
	instance := instance.New(appName)
	device := instance.GetDevice(deviceIndex)

	// transfer worker
	transfer := command.NewWorker(device, vk.QueueFlags(vk.QueueTransferBit), 0)

	// per frame graphics workers
	workerCount := 1 // frames
	workers := make([]command.Worker, workerCount)
	for i := range workers {
		workers[i] = command.NewWorker(device, vk.QueueFlags(vk.QueueGraphicsBit), i+1)
	}

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
	return b.workers[frame%len(b.workers)]
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

	// destroy workers
	b.transfer.Destroy()
	for _, w := range b.workers {
		w.Destroy()
	}
	b.workers = nil

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
