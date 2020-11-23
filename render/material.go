package render

import (
	"fmt"
)

// Material contains a shader reference and all resources required to draw a vertex buffer array
type Material struct {
	*Shader
	Textures *TextureMap

	name string
}

// CreateMaterial instantiates a new empty material
func CreateMaterial(name string, shader *Shader) *Material {
	return &Material{
		Shader:   shader,
		Textures: NewTextureMap(shader),
		name:     name,
	}
}

func (mat *Material) String() string {
	return fmt.Sprintf("Material %s", mat.name)
}

// Use sets the current shader and activates textures
func (mat *Material) Use() {
	mat.Shader.Use()
	mat.Textures.Use()
}
