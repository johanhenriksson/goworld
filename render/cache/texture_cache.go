package cache

import (
	"github.com/johanhenriksson/goworld/render/command"
	"github.com/johanhenriksson/goworld/render/device"
	"github.com/johanhenriksson/goworld/render/texture"
	"github.com/johanhenriksson/goworld/render/upload"
)

type TextureCache T[texture.Ref, texture.T]

// mesh cache backend
type textures struct {
	device device.T
	worker command.Worker
}

func NewTextureCache(dev device.T, transferer command.Worker) TextureCache {
	return New[texture.Ref, texture.T](&textures{
		device: dev,
		worker: transferer,
	})
}

func (t *textures) Name() string {
	return "Texture"
}

func (t *textures) Instantiate(ref texture.Ref) texture.T {
	img := ref.Load()

	tex, err := upload.NewTextureSync(t.device, t.worker, img)
	if err != nil {
		panic(err)
	}

	return tex
}

func (m *textures) Update(tex texture.T, ref texture.Ref) texture.T {
	// we cant reuse texture objects yet
	tex2 := m.Instantiate(ref)
	tex.Destroy()
	return tex2
}

func (m *textures) Delete(tex texture.T) {
	tex.Destroy()
}

func (m *textures) Destroy() {
}