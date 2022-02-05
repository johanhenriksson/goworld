package material

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
)

type T interface {
	shader.T

	Texture(name string, tex texture.T)
	TextureSlot(slot int) texture.T
}

// material contains a shader reference and all resources required to draw a vertex buffer array
type material struct {
	shader.T

	name     string
	textures map[string]texture.T
	slots    []string
}

// New instantiates a new empty material
func New(name string, shader shader.T) T {
	return &material{
		T:        shader,
		name:     name,
		textures: make(map[string]texture.T, 8),
		slots:    make([]string, 0, 8),
	}
}

func (mat *material) String() string {
	return fmt.Sprintf("Material %s", mat.name)
}

// Use sets the current shader and activates textures
func (mat *material) Use() error {
	if err := mat.T.Use(); err != nil {
		return err
	}
	for i, name := range mat.slots {
		slot := texture.Slot(i)
		tex := mat.textures[name]
		if err := tex.Use(slot); err != nil {
			return err
		}
		if err := mat.T.Texture2D(name, slot); err != nil {
			return err
		}
	}
	return nil
}

func (mat *material) Texture(name string, tex texture.T) {
	if _, exists := mat.textures[name]; exists {
		mat.textures[name] = tex
		return
	}
	mat.textures[name] = tex
	mat.slots = append(mat.slots, name)
}

func (mat *material) TextureSlot(slot int) texture.T {
	if slot < 0 || slot >= len(mat.slots) {
		return nil
	}
	return mat.textures[mat.slots[slot]]
}
