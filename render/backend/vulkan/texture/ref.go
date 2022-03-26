package texture

import (
	"image"

	"github.com/johanhenriksson/goworld/engine/cache"
	"github.com/johanhenriksson/goworld/render"
)

type Ref interface {
	cache.Item
	Load() *image.RGBA
}

type path_ref struct {
	path string
}

func PathRef(path string) Ref {
	return &path_ref{path}
}

func (r *path_ref) Id() string   { return r.path }
func (r *path_ref) Version() int { return 1 }

func (r *path_ref) Load() *image.RGBA {
	img, err := render.ImageFromFile(r.path)
	if err != nil {
		panic(err)
	}
	return img
}

type image_ref struct {
	name string
	img  *image.RGBA
}

func ImageRef(name string, img *image.RGBA) Ref {
	return &image_ref{
		name: name,
		img:  img,
	}
}

func (r *image_ref) Id() string        { return r.name }
func (r *image_ref) Version() int      { return 1 }
func (r *image_ref) Load() *image.RGBA { return r.img }
