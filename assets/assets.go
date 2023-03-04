package assets

import (
	"fmt"
	"image"
	"image/draw"
	"io/fs"
	"os"
	"path/filepath"

	// image codecs
	_ "image/png"

	"github.com/johanhenriksson/goworld/render/font"
)

type FontMap map[string]font.T

type ResourceCache struct {
	Fonts FontMap
}

/* Global asset cache */
var cache *ResourceCache
var vfs fs.FS

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	assetRoot := FindFileInParents("assets", cwd)
	vfs = os.DirFS(assetRoot)

	cache = &ResourceCache{
		Fonts: make(FontMap),
	}
}

func Open(file string) (fs.File, error) {
	return vfs.Open(file)
}

func GetImage(file string) (*image.RGBA, error) {
	imgFile, err := Open(file)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	return rgba, nil
}

func GetFont(name string, size int, scale float32) font.T {
	key := fmt.Sprintf("%s-%d", name, size)
	if font, exists := cache.Fonts[key]; exists {
		return font
	}

	fmt.Printf("+ font %s %dpt\n", name, size)

	file, err := vfs.Open(name)
	if err != nil {
		panic(fmt.Errorf("error reading font %s: %w", name, err))
	}

	font, err := font.Load(file, int(float32(size)*scale))
	if err != nil {
		panic(fmt.Errorf("error loading font %s: %w", name, err))
	}

	cache.Fonts[key] = font
	return font
}

func FindFileInParents(name, path string) string {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.Name() == name {
			return filepath.Join(path, name)
		}
	}
	return FindFileInParents(name, filepath.Dir(path))
}
