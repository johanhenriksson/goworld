package engine

import (
	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/descriptor"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/instance"
)

type App interface {
	Instance() *instance.Instance
	Device() *device.Device
	Destroy()

	Worker() command.Worker
	Flush()

	Pool() *descriptor.Pool
	Meshes() cache.MeshCache
	Textures() cache.TextureCache
	Shaders() cache.ShaderCache
}

type engine struct {
	appName  string
	instance *instance.Instance
	device   *device.Device

	worker   command.Worker
	pool     *descriptor.Pool
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

	// workers
	worker := command.NewWorker(device, "worker", device.Queue())

	// init caches
	meshes := cache.NewMeshBlockCache(device, worker, 1024*1024*1024, 256*1024*1024)
	textures := cache.NewTextureCache(device, worker)
	shaders := cache.NewShaderCache(device)

	pool := descriptor.NewPool(device, 1000, DefaultDescriptorPools)

	return &engine{
		appName: appName,

		device:   device,
		instance: instance,
		worker:   worker,
		meshes:   meshes,
		textures: textures,
		shaders:  shaders,
		pool:     pool,
	}
}

func (b *engine) Instance() *instance.Instance { return b.instance }
func (b *engine) Device() *device.Device       { return b.device }

func (b *engine) Pool() *descriptor.Pool       { return b.pool }
func (b *engine) Meshes() cache.MeshCache      { return b.meshes }
func (b *engine) Textures() cache.TextureCache { return b.textures }
func (b *engine) Shaders() cache.ShaderCache   { return b.shaders }

func (b *engine) Worker() command.Worker {
	return b.worker
}

func (b *engine) Flush() {
	// wait for all workers
	b.worker.Flush()
	b.device.WaitIdle()
}

func (b *engine) Destroy() {
	// flush any pending work
	b.Flush()

	// clean up caches
	b.pool.Destroy()
	b.meshes.Destroy()
	b.textures.Destroy()
	b.shaders.Destroy()

	// destroy workers
	b.worker.Destroy()

	if b.device != nil {
		b.device.Destroy()
		b.device = nil
	}
	if b.instance != nil {
		b.instance.Destroy()
		b.instance = nil
	}
}
