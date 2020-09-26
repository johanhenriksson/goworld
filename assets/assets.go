package assets

import (
	"fmt"

	"github.com/johanhenriksson/goworld/render"
)

type ShaderMap map[string]*render.ShaderProgram
type TextureMap map[string]*render.Texture
type MaterialMap map[string]*render.Material
type FontMap map[string]*render.Font

type ResourceCache struct {
	Shaders   ShaderMap
	Textures  TextureMap
	Materials MaterialMap
	Fonts     FontMap
}

/* Global asset cache */
var cache *ResourceCache

func init() {
	cache = &ResourceCache{
		Shaders:   make(ShaderMap),
		Textures:  make(TextureMap),
		Materials: make(MaterialMap),
		Fonts:     make(FontMap),
	}
}

func GetShader(name string) *render.ShaderProgram {
	// check shader cache
	if shader, exists := cache.Shaders[name]; exists {
		return shader
	}

	// attempt to load
	fmt.Println("+ shader", name)
	shader := render.CompileShaderProgram("assets/shaders/" + name)
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
	texture, error := render.TextureFromFile("assets/" + name)
	if error != nil {
		panic(fmt.Sprintf("Error loading texture %s: %s", name, error))
	}

	cache.Textures[name] = texture
	return texture
}

func GetFont(name string, size, spacing float32) *render.Font {
	key := fmt.Sprintf("%s-%.1f-%.1f", name, size, spacing)
	if font, exists := cache.Fonts[key]; exists {
		return font
	}

	dpi := float32(1.0)
	fmt.Printf("+ font %s (%.1fpt, %.1f, %dx)\n", name, size, spacing, int(dpi))
	font := render.LoadFont(name, dpi, size, spacing)
	cache.Fonts[key] = font

	return font
}
