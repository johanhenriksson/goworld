package assets

import (
	"fmt"
	"github.com/johanhenriksson/goworld/render"
)

type ShaderMap map[string]*render.ShaderProgram
type TextureMap map[string]*render.Texture
type MaterialMap map[string]*render.Material

type ResourceCache struct {
	Shaders   ShaderMap
	Textures  TextureMap
	Materials MaterialMap
}

/* Global asset cache */
var cache *ResourceCache

func init() {
	cache = &ResourceCache{
		Shaders:   make(ShaderMap),
		Textures:  make(TextureMap),
		Materials: make(MaterialMap),
	}
}

func GetShader(name string) *render.ShaderProgram {
	// check shader cache
	if shader, exists := cache.Shaders[name]; exists {
		return shader
	}

	// attempt to load
	fmt.Println("+ shader", name)
	shader := render.CompileVFShader("assets/shaders/" + name)
	cache.Shaders[name] = shader

	return shader
}

func GetTexture(name string) *render.Texture {
	// check texture cache
	if texture, exists := cache.Textures[name]; exists {
		return texture
	}

	// attempt to load
	fmt.Println("+ texture", name)
	texture, error := render.LoadTexture("assets/" + name)
	if error != nil {
		panic(fmt.Sprintf("Error loading texture %s: %s", name, error))
	}

	cache.Textures[name] = texture
	return texture
}

func GetMaterial(name string) *render.Material {
	/* load shader by the same name */
	shader := GetShader(name)

	// attempt to load
	fmt.Println("+ material", name)
	material := render.LoadMaterial(shader, "assets/materials/"+name)

	return material
}

func GetMaterialCached(name string) *render.Material {
	if mat, exists := cache.Materials[name]; exists {
		return mat
	}

	mat := GetMaterial(name)
	cache.Materials[name] = mat

	return mat
}
