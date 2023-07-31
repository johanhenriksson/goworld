package vulkan

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vulkan/instance"

	"github.com/vkngwrapper/core/v2/core1_0"
)

type App interface {
	Instance() instance.T
	Device() device.T
	Destroy()

	Worker(int) command.Worker
	Transferer() command.Worker
	Flush()

	Pool() descriptor.Pool
	Meshes() cache.MeshCache
	Textures() cache.TextureCache
	Shaders() cache.ShaderCache
}

type backend struct {
	appName  string
	instance instance.T
	device   device.T

	transfer command.Worker
	workers  []command.Worker

	pool     descriptor.Pool
	meshes   cache.MeshCache
	textures cache.TextureCache
	shaders  cache.ShaderCache
}

func New(appName string, deviceIndex int) App {
	instance := instance.New(appName)
	device, err := device.New(instance, instance.EnumeratePhysicalDevices()[0])
	if err != nil {
		panic(err)
	}

	// transfer worker
	transfer := command.NewWorker(device, "xfer", core1_0.QueueTransfer|core1_0.QueueGraphics, 0)

	// per frame graphics workers
	workerCount := 1 // frames
	workers := make([]command.Worker, workerCount)
	for i := range workers {
		workers[i] = command.NewWorker(device, fmt.Sprintf("frame%d", i), core1_0.QueueGraphics, i+1)
	}

	// init caches
	meshes := cache.NewMeshCache(device, transfer)
	textures := cache.NewTextureCache(device, transfer)
	shaders := cache.NewShaderCache(device)

	pool := descriptor.NewPool(device, DefaultDescriptorPools)

	return &backend{
		appName: appName,

		device:   device,
		instance: instance,
		transfer: transfer,
		workers:  workers,
		meshes:   meshes,
		textures: textures,
		shaders:  shaders,
		pool:     pool,
	}
}

func (b *backend) Instance() instance.T { return b.instance }
func (b *backend) Device() device.T     { return b.device }

func (b *backend) Pool() descriptor.Pool        { return b.pool }
func (b *backend) Meshes() cache.MeshCache      { return b.meshes }
func (b *backend) Textures() cache.TextureCache { return b.textures }
func (b *backend) Shaders() cache.ShaderCache   { return b.shaders }

func (b *backend) Transferer() command.Worker {
	return b.transfer
}

func (b *backend) Worker(frame int) command.Worker {
	return b.workers[frame%len(b.workers)]
}

func (b *backend) Flush() {
	// wait for all workers
	for _, w := range b.workers {
		w.Flush()
	}
	b.device.WaitIdle()
}

func (b *backend) Destroy() {
	// flush any pending work
	b.Flush()

	// clean up caches
	b.pool.Destroy()
	b.meshes.Destroy()
	b.textures.Destroy()
	b.shaders.Destroy()

	// destroy workers
	b.transfer.Destroy()
	for _, w := range b.workers {
		w.Destroy()
	}
	b.workers = nil

	if b.device != nil {
		b.device.Destroy()
		b.device = nil
	}
	if b.instance != nil {
		b.instance.Destroy()
		b.instance = nil
	}
}
