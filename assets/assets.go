package assets

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/johanhenriksson/goworld/core/window"
	"github.com/johanhenriksson/goworld/render/font"
)

type FontMap map[string]font.T

type ResourceCache struct {
	Fonts FontMap
}

/* Global asset cache */
var cache *ResourceCache

func init() {
	cache = &ResourceCache{
		Fonts: make(FontMap),
	}
}

func GetFont(name string, size int) font.T {
	key := fmt.Sprintf("%s-%d", name, size)
	if font, exists := cache.Fonts[key]; exists {
		return font
	}

	fmt.Printf("+ font %s %dpt\n", name, size)
	font := font.Load(AssetPath(name), int(float32(size)*window.Scale))
	cache.Fonts[key] = font

	return font
}

func DefaultFont() font.T {
	return GetFont("fonts/SourceCodeProRegular.ttf", 12)
}

var assetRoot = ""

func AssetPath(path string, args ...any) string {
	if assetRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		assetRoot = FindFileInParents("assets", cwd)
	}
	return filepath.Join(assetRoot, "assets", fmt.Sprintf(path, args...))
}

func FindFileInParents(name, path string) string {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.Name() == name {
			return path
		}
	}
	return FindFileInParents(name, filepath.Dir(path))
}
