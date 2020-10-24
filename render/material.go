package render

import (
	"fmt"
)

// MaterialTextureMap maps texture names to texture objects
type MaterialTextureMap map[string]*Texture

// Material contains a shader reference and all resources required to draw a vertex buffer array
type Material struct {
	*Shader
	Textures MaterialTextureMap

	texslots []string // since map is unordered
	name     string
}

// CreateMaterial instantiates a new empty material
func CreateMaterial(name string, shader *Shader) *Material {
	return &Material{
		Shader:   shader,
		Textures: make(MaterialTextureMap),
		name:     name,
	}
}

func (mat *Material) String() string {
	return fmt.Sprintf("Material %s", mat.name)
}

// AddTexture attaches a new texture to this material, and assings it to the next available texture slot.
func (mat *Material) AddTexture(name string, tex *Texture) {
	if _, exists := mat.Textures[name]; exists {
		mat.SetTexture(name, tex)
		return
	}
	mat.Textures[name] = tex
	mat.texslots = append(mat.texslots, name)
}

// SetTexture changes a bound texture
func (mat *Material) SetTexture(name string, tex *Texture) {
	mat.Textures[name] = tex
}

// Use sets the current shader and activates textures
func (mat *Material) Use() {
	mat.Shader.Use()
	i := uint32(0)
	for _, name := range mat.texslots {
		tex := mat.Textures[name]
		tex.Use(i)
		mat.Int32(name, int32(i))
		i++
	}
}
