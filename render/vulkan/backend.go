package vulkan

import (
	"github.com/johanhenriksson/goworld/render/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/vulkan/instance"
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
	graphics command.Worker

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

	// general purpose worker
	worker := command.NewWorker(device, "all", device.GraphicsQueue())
	transfer := worker
	graphics := worker

	// transfer := command.NewWorker(device, "xfer", device.TransferQueue())
	// graphics := command.NewWorker(device, "graphics", device.GraphicsQueue())

	// init caches
	meshes := cache.NewMeshCache(device, transfer)
	textures := cache.NewTextureCache(device, transfer)
	shaders := cache.NewShaderCache(device)

	pool := descriptor.NewPool(device, 1000, DefaultDescriptorPools)

	return &backend{
		appName: appName,

		device:   device,
		instance: instance,
		transfer: transfer,
		graphics: graphics,
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
	return b.graphics
}

func (b *backend) Flush() {
	// wait for all workers
	b.transfer.Flush()
	b.graphics.Flush()
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
	b.graphics.Destroy()
	b.graphics = nil

	if b.device != nil {
		b.device.Destroy()
		b.device = nil
	}
	if b.instance != nil {
		b.instance.Destroy()
		b.instance = nil
	}
}
