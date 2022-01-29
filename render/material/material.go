package material

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
)

type T interface {
	shader.T

	GetTexture(name string) texture.T
	SetTexture(name string, tex texture.T)
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
		textures: make(map[string]texture.T),
		slots:    make([]string, 0, 8),
	}
}

func (mat *material) String() string {
	return fmt.Sprintf("Material %s", mat.name)
}

// Use sets the current shader and activates textures
func (mat *material) Use() {
	mat.T.Use()
	for i, name := range mat.slots {
		tex := mat.textures[name]
		tex.Use(i)
		mat.T.Int32(name, i)
	}
}

func (mat *material) GetTexture(name string) texture.T {
	if tx, exists := mat.textures[name]; exists {
		return tx
	}
	return nil
}

func (mat *material) SetTexture(name string, tex texture.T) {
	if _, exists := mat.textures[name]; exists {
		mat.textures[name] = tex
		return
	}
	mat.textures[name] = tex
	mat.slots = append(mat.slots, name)
}
