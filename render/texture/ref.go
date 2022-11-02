package texture

import (
	"image"

	"github.com/johanhenriksson/goworld/assets"
)

type Ref interface {
	Id() string
	Version() int

	// Load is called by texture caches and loaders, and should return the image data.
	// todo: This interface is a bit too simple as it does not allow us to pass
	//       formats, filters and aspects.
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
	img, err := assets.GetImage(r.path)
	if err != nil {
		panic(err)
	}
	return img
}

type image_ref struct {
	name    string
	version int
	img     *image.RGBA
}

func ImageRef(name string, version int, img *image.RGBA) Ref {
	return &image_ref{
		name:    name,
		version: version,
		img:     img,
	}
}

func (r *image_ref) Id() string        { return r.name }
func (r *image_ref) Version() int      { return r.version }
func (r *image_ref) Load() *image.RGBA { return r.img }
