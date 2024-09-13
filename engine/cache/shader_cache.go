package cache

import (
	"github.com/johanhenriksson/goworld/assets"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/shader"
)

type ShaderCache T[assets.Shader, *shader.Shader]

func NewShaderCache(dev *device.Device) ShaderCache {
	return New[assets.Shader, *shader.Shader](&shaders{
		device: dev,
	})
}

type shaders struct {
	device *device.Device
}

func (s *shaders) Name() string {
	return "Shaders"
}

func (s *shaders) Instantiate(key assets.Shader, callback func(*shader.Shader)) {
	// load shader in a background goroutine
	go func() {
		shader := key.LoadShader(assets.FS, s.device)
		callback(shader)
	}()
}

func (s *shaders) Delete(shader *shader.Shader) {
	shader.Destroy()
}

func (s *shaders) Destroy() {
}
