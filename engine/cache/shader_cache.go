package cache

import (
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/shader"
)

type ShaderCache T[shader.Ref, shader.T]

func NewShaderCache(dev device.T) ShaderCache {
	return New[shader.Ref, shader.T](&shaders{
		device: dev,
	})
}

type shaders struct {
	device device.T
}

func (s *shaders) Name() string {
	return "Shaders"
}

func (s *shaders) Instantiate(key shader.Ref, callback func(shader.T)) {
	// load shader in a background goroutine
	go func() {
		shader := key.Load(s.device)
		callback(shader)
	}()
}

func (s *shaders) Delete(shader shader.T) {
	shader.Destroy()
}

func (s *shaders) Destroy() {
}
