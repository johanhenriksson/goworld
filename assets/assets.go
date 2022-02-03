package assets

import (
	"fmt"
	"os"

	"github.com/johanhenriksson/goworld/render"
	glshader "github.com/johanhenriksson/goworld/render/backend/gl/gl_shader"
	"github.com/johanhenriksson/goworld/render/color"
	"github.com/johanhenriksson/goworld/render/font"
	"github.com/johanhenriksson/goworld/render/material"
	"github.com/johanhenriksson/goworld/render/shader"
	"github.com/johanhenriksson/goworld/render/texture"
)

type ShaderMap map[string]shader.T
type TextureMap map[string]texture.T
type MaterialMap map[string]material.T
type FontMap map[string]font.T

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

func GetShader(name string) shader.T {
	// check shader cache
	if shader, exists := cache.Shaders[name]; exists {
		return shader
	}

	// attempt to load
	fmt.Println("+ shader", name)

	files := []string{
		fmt.Sprintf("assets/shaders/%s.vs", name),
		fmt.Sprintf("assets/shaders/%s.fs", name),
	}

	// optional geometry shader
	gsPath := fmt.Sprintf("assets/shaders/%s.gs", name)
	if _, err := os.Stat(gsPath); err == nil {
		files = append(files, gsPath)
	}

	shader := glshader.CompileShader(name, files...)
	cache.Shaders[name] = shader

	return shader
}

func GetTexture(name string) texture.T {
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

func GetColorTexture(color color.T) texture.T {
	name := fmt.Sprintf("Color%s", color)

	// check texture cache
	if texture, exists := cache.Textures[name]; exists {
		return texture
	}

	// otherwise, instantiate a new one
	fmt.Println("+ texture", name)
	texture := render.TextureFromColor(color)
	cache.Textures[name] = texture
	return texture
}

func GetFont(name string, size int) font.T {
	key := fmt.Sprintf("%s-%d", name, size)
	if font, exists := cache.Fonts[key]; exists {
		return font
	}

	fmt.Printf("+ font %s %dpt\n", name, size)
	font := font.Load(name, size)
	cache.Fonts[key] = font

	return font
}
