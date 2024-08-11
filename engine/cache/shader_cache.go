package cache

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/shader"
)

type ShaderCache T[shader.Ref, *shader.Shader]

func NewShaderCache(dev *device.Device) ShaderCache {
	return New[shader.Ref, *shader.Shader](&shaders{
		device: dev,
	})
}

type shaders struct {
	device *device.Device
}

func (s *shaders) Name() string {
	return "Shaders"
}

func (s *shaders) Instantiate(key shader.Ref, callback func(*shader.Shader)) {
	// load shader in a background goroutine
	go func() {
		shader := key.Load(s.device)
		callback(shader)
	}()
}

func (s *shaders) Delete(shader *shader.Shader) {
	shader.Destroy()
}

func (s *shaders) Destroy() {
}
